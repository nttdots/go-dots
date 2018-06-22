package data_messages

import (
  "encoding/json"
  "fmt"
  "net/http"
  "reflect"

  "github.com/fatih/structs"
  log "github.com/sirupsen/logrus"
)

type Content string

const (
  Content_Config    Content = "config"
  Content_NonConfig Content = "nonconfig"
  Content_All       Content = "all"
)

func ContentFromRequest(r *http.Request) (Content, error) {
  q := r.URL.Query()
  if a, ok := q["content"]; ok {
    if len(a) == 1 {
      switch a[0] {
      case string(Content_Config):
        return Content_Config, nil
      case string(Content_NonConfig):
        return Content_NonConfig, nil
      case string(Content_All):
        return Content_All, nil
      default:
        return Content_All, fmt.Errorf("Unknown content parameter: %s", a[0])
      }
    } else {
      return Content_All, fmt.Errorf("Multiple 'content' parameter specified: len=%d", len(a))
    }
  } else {
    return Content_All, nil
  }
}

/*
 * recursive struct to map[string]interface{} conversion.
 * * use json tag for field name.
 * * use yang tag for filtering.
 *                     `yang:"config"` `yang:"nonconfig"` -
 *   Content_Config    yes             -                  yes
 *   Content_NonConfig -               yes                yes
 *   Content_All       yes             yes                yes

 */
func convert(v interface{}, content Content) (interface{}, error) {
  if v == nil {
    return nil, nil
  }

  value := reflect.ValueOf(v)

  switch value.Kind() {
  case reflect.Chan, reflect.Func, reflect.Interface:
    if value.IsNil() {
      return nil, nil
    }
    return v, nil

  case reflect.Map:
    if value.IsNil() {
      return nil, nil
    }
    if isJsonMarshaler(v) {
      return v, nil
    } else {
      r := make(map[string]interface{})
      if value.Len() == 0 { // preserve empty map
        return r, nil
      }

      for _, k := range value.MapKeys() {
        cv, err := convert(value.MapIndex(k).Interface(), content)
        if err != nil {
          return nil, err
        }
        if cv != nil {
          r[fmt.Sprintf("%v", k)] = cv // use %v anyway
        }
      }
      if len(r) == 0 {
        return nil, nil
      }
      return r, nil
    }

  case reflect.Slice:
    if value.IsNil() {
      return nil, nil
    }
    fallthrough
  case reflect.Array:
    if isJsonMarshaler(v) {
      return v, nil
    } else {
      r := make([]interface{}, 0)
      if value.Len() == 0 { // preserve empty slice
        return r, nil
      }

      for i := 0; i < value.Len(); i++ {
        cv, err := convert(value.Index(i).Interface(), content)
        if err != nil {
          return nil, err
        }
        if cv != nil {
          r = append(r, cv) // do not filter nil
        }
      }
      if len(r) == 0 {
        return nil, nil
      }
      return r, nil
    }

  case reflect.Ptr:
    if value.IsNil() {
      return nil, nil
    }
    return convert(value.Elem().Interface(), content)

  case reflect.Struct:
    if isJsonMarshaler(v) {
      return v, nil
    } else {
      r := make(map[string]interface{})
      for _, f := range structs.New(v).Fields() {
        if !f.IsExported() {
          continue
        }

        yang := f.Tag("yang")
        switch yang {
        case "":
          // no action
        case "config":
          if content == Content_NonConfig {
            continue
          }
        case "nonconfig":
          if content == Content_Config {
            continue
          }
        default:
          log.WithField("yang", yang).Warn("Unexpected `yang` tag.")
          continue
        }

        name := f.Tag("json")
        if name == "" {
          name = f.Name()
        }

        fv := f.Value()
        cv, err := convert(fv, content)
        if err != nil {
          return nil, err
        }

        if cv != nil {
          r[name] = cv
        }
      }
      if len(r) == 0 {
        return nil, nil
      }
      return r, nil
    }

  default:
    return v, nil
  }
}

func ToMap(src interface{}, content Content) (map[string]interface{}, error) {
  r, err := convert(src, content)
  if err != nil {
    return nil, err
  }

  if r == nil {
    return make(map[string]interface{}), nil
  }

  if m, ok := r.(map[string]interface{}); ok {
    return m, nil
  } else {
    return nil, fmt.Errorf("Converted result is not map.")
  }
}

var marshalerType = reflect.TypeOf(new(json.Marshaler)).Elem()

func isJsonMarshaler(x interface{}) bool {
  t := reflect.TypeOf(x)

  if t.Implements(marshalerType) {
    return true
  }

  if reflect.PtrTo(t).Implements(marshalerType) {
    return true
  }

  return false
}
