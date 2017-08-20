package txt

func GetLine( str string ) (string, int) {
  for i, c := range str {
    if c == '\n' {
      return str[:i], i + 1
    }
  }

  return str, len( str )
}

func GetLines( str string ) []string {
  result := make( []string, 0, 64 )
  last := 0;
  for i, c := range str {
    if c == '\n' {
      result = append( result, str[last:i] )
      last = i + 1
    }
  }

  if last < len( str ) { result = append( result, str[last:] ) }

  return result
}

func RmSpacesAtEnd( str string ) string {
  for i := len( str ) - 1; i >= 0; i-- {
    switch str[i] {
    case ' ', '\t', '\n', '\v', '\f', '\r' :
    default: return str[:i+1]
    }
  }

  return ""
}

func RmSpacesAtStartup( str string ) string {
  for i := 0; i < len( str ); i++ {
    switch str[i] {
    case ' ', '\t', '\n', '\v', '\f', '\r' :
    default: return str[i:]
    }
  }

  return ""
}

func RmSpacesToTheSides( str string ) string {
  return RmSpacesAtEnd( RmSpacesAtStartup( str ) )
}

func Linelize( str string ) string {
  lines := GetLines( str )
  if len( lines ) == 0 { return "" }

  k   := make( []byte, len( str ) )
  pos := copy( k, RmSpacesToTheSides( lines[0] ) )

  for i := 1; i < len( lines ); i++ {
    if len( lines[ i ] ) == 0 { continue }
    if pos > 0 && k[pos - 1 ] != ' ' {
      k[pos] = ' ';
      pos++
    }
    pos += copy( k[pos:], RmSpacesToTheSides( lines[i] ) )
  }

  return string( k[:pos] )
}

func SpaceSwap( str, swap string ) string {
  i, j, k := 0, 0, make( []byte, len( str ) + len( swap ) * countSpacesRegions( str ) )

  for ; i < len( str );  {
    switch str[i] {
    case ' ', '\t', '\n', '\v', '\f', '\r' :
      i += countInitSpaces( str[i:] )
      copy( k[j:], swap[0:] )
      j += len( swap )
    default: k[j] = str[i]; j++; i++
    }
  }

  return string( k[:j] )
}

func countSpacesRegions( str string ) (n int) {
  for i := 0; i < len( str ); i++ {
    switch str[i] {
    case ' ', '\t', '\n', '\v', '\f', '\r' :
      i += countInitSpaces( str[i:] )
      n++
    }
  }

  return
}

func RmIndent( str string, indentLevel int ) string {
  i, j, k := 0, 0, make( []byte, len( str ) )

  if CountIndentSpaces( str ) >= indentLevel {
    i = indentLevel
  }

  for ; i < len( str ); i++ {
    k[j] = str[i]
    j++

    if str[ i ] == '\n' {
      if CountIndentSpaces( str[ i + 1:] ) >= indentLevel {
        i += indentLevel
      }
    }
  }

  return string( k[:j] )
}

func CountIndentSpaces( str string ) int {
  for i := 0; i < len( str ); i++ {
    switch str[i] {
    case ' ', '\t':
    default: return i
    }
  }

  return len(str)
}

func DragTextByIndent( str string, indent int ) (string, int) {
  for init, width, line := 0, 0, ""; init < len(str); {
    line, width = GetLine( str[init:] )

    if countInitSpaces( line ) < indent {
      return str[:init], init
    }

    init += width
  }

  return str, len(str)
}

func DragAllTextByIndent( str string, indent int ) (string, int) {
  for init, width, line := 0, 0, ""; init < len(str); {
    line, width = GetLine( str[init:] )

    if countInitSpaces( line ) >= indent || len(line) == 0 {
      init += width
      continue
    }

    return str[:init], init
  }

  return str, len(str)
}

func countInitChars( str string ) int {
  for i, c := range str {
    switch c {
    case ' ', '\t', '\n', '\v', '\f', '\r' : return i
    }
  }

  return len( str )
}

func countInitSpaces( str string ) int {
  for i, c := range str {
    switch c {
    case ' ', '\t', '\n', '\v', '\f', '\r' :
    default: return i
    }
  }

  return len(str)
}

func Tokenize( str string ) []string {
  r := make([]string, 0, 16)

  i, w, max := 0, 0, len( str )
  for i < max {
    i += countInitSpaces( str[i:] )
    w  = countInitChars ( str[i:] )

    if w > 0 {
      r = append( r, str[i:i+w] )
      i += w
    } else { break }
  }

  return r
}
