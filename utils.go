package goopt_fluent

import (

)

func quoteJoin(word string, length int, get func (int) string) (str string) {
   str = "'" + get(0) + "'"
   var i int
   for i = 1; i + 1 < length; i++ {
      str += ", '" + get(i) + "'"
   }
   str += word + "'" + get(i) + "'"
   return
}