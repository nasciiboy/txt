package txt

import (
  "regexp"
  "testing"
  "bytes"

  "github.com/nasciiboy/regexp4"
)

func TestGetLine( t *testing.T ){
  data := []struct{
    input    string
    output   string
    n int
  } {
    { "", "", 0 },
    { "\n\n\n", "", 1 },
    { "line", "line", 4 },
    { "line\n", "line", 5 },
    { "line\t\v\rline", "line\t\v\rline", 11 },
    { "line\t\v\nline\n", "line\t\v", 7 },
    { "line\t\v\n\nline\n\n", "line\t\v", 7 },
    { "line\t\v\n\n\tline\n\n\n", "line\t\v", 7 },
    { "\n1\n2\n3\n", "", 1 },
    { "1\n2\n3\n4\n5\n", "1", 2 },
  }

  for _, d := range data {
    output, dout := GetLine( d.input )
    if output != d.output || dout != d.n {
      t.Errorf( "GetLine( %q ) \nreturn   [%d] %q\nexpected [%d] %q", d.input, dout, output, d.n, d.output )
    }
  }
}

func TestGetRawLine( t *testing.T ){
  data := []struct{
    input    string
    output   string
  } {
    { "", "" },
    { "\n\n\n", "\n" },
    { "line", "line" },
    { "line\n", "line\n" },
    { "line\t\v\rline", "line\t\v\rline" },
    { "line\t\v\nline\n", "line\t\v\n" },
    { "line\t\v\n\nline\n\n", "line\t\v\n" },
    { "line\t\v\n\n\tline\n\n\n", "line\t\v\n" },
    { "\n1\n2\n3\n", "\n" },
    { "1\n2\n3\n4\n5\n", "1\n" },
  }

  for _, d := range data {
    output := GetRawLine( d.input )
    if output != d.output {
      t.Errorf( "GetRawLine( %q ) \nreturn   %q\nexpected %q", d.input, output, d.output )
    }
  }
}

func TestGetLines( t *testing.T ){
  data := []struct{
    input    string
    output   []string
    lines    int
  } {
    { "", []string{}, 0 },
    { "\n", []string{ "" }, 1 },
    { "\n\n\n", []string{ "", "", "" }, 3 },
    { "\n\ta", []string{ "", "\ta" }, 2 },
    { "line", []string{ "line" }, 1 },
    { "line\t\v\rline", []string{ "line\t\v\rline" }, 1 },
    { "line\t\v\rline\n", []string{ "line\t\v\rline" }, 1 },
    { "line\t\v\rline\n\n", []string{ "line\t\v\rline", "", }, 2 },
    { "line\t\v\rline\n\n\n", []string{ "line\t\v\rline", "", "" }, 3 },
    { "\n1\n2\n3\n", []string{ "", "1", "2", "3" }, 4 },
    { "1\n2\n3\n4\n5\n", []string{ "1", "2", "3", "4", "5" }, 5 },
  }

  for _, d := range data {
    output := GetLines( d.input )
    if !cmpStringArray( output, d.output ) || len( output ) != d.lines {
      t.Errorf( "GetLines( %q ) \nreturn   [%d] %q\nexpected [%d] %q", d.input, len(output), output, d.lines, d.output )
    }
  }
}

func TestGetRawLines( t *testing.T ){
  data := []struct{
    input    string
    output   []string
    lines    int
  } {
    { "", []string{}, 0 },
    { "\n", []string{ "\n" }, 1 },
    { "\n\n\n", []string{ "\n", "\n", "\n" }, 3 },
    { "\n\ta", []string{ "\n", "\ta" }, 2 },
    { "line", []string{ "line" }, 1 },
    { "line\t\v\rline", []string{ "line\t\v\rline" }, 1 },
    { "line\t\v\rline\n", []string{ "line\t\v\rline\n" }, 1 },
    { "line\t\v\rline\n\n", []string{ "line\t\v\rline\n", "\n", }, 2 },
    { "line\t\v\rline\n\n\n", []string{ "line\t\v\rline\n", "\n", "\n" }, 3 },
    { "\n1\n2\n3\n", []string{ "\n", "1\n", "2\n", "3\n" }, 4 },
    { "1\n2\n3\n4\n5\n", []string{ "1\n", "2\n", "3\n", "4\n", "5\n" }, 5 },
  }

  for _, d := range data {
    output := GetRawLines( d.input )
    if !cmpStringArray( output, d.output ) || len( output ) != d.lines {
      t.Errorf( "GetLines( %q ) \nreturn   [%d] %q\nexpected [%d] %q", d.input, len(output), output, d.lines, d.output )
    }
  }
}

func TestTokenize( t *testing.T ){
  data := []struct{
    input    string
    output   []string
    expected bool
  } {
    { "", nil, true },
    { "a b c", nil, false },
    { "a b c", []string{ "a", "b", "c" }, true },
    { "a b c", []string{ "a", "b", "d" }, false },
    { "hola, que tal!", []string{ "hola,", "que", "tal!" }, true },
    { "hola,\n que\n tal!", []string{ "hola,", "que", "tal!" }, true },
    { "hola,\n\n\n\nque\t\v tal!", []string{ "hola,", "que", "tal!" }, true },
    { " \n\t  hola,\n\n\n\n que\t\v tal!", []string{ "hola,", "que", "tal!" }, true },
    { "–bueno–, que es esto? nada...", []string{ "–bueno–,", "que", "es", "esto?", "nada..." }, true },
    { "como–es esto–que es esto?", []string{ "como–es", "esto–que", "es", "esto?" }, true },
    { "como (es esto). Que es esto", []string{ "como", "(es", "esto).", "Que", "es", "esto" }, true },
    { "como (es esto). Que es esto", []string{ "como", "(es", "esto)." }, false },
  }

  for _, d := range data {
    output := Tokenize( d.input )
    if cmpStringArray( output, d.output ) != d.expected {
      t.Errorf( "Tokenize( %q ) \nreturn   %v\nexpected %v", d.input, output, d.output )
    }
  }
}

func cmpStringArray( a, b []string ) bool {
  if len( a ) != len( b ) { return false }

  for i, s := range a {
    if s != b[i] { return false }
  }

  return true
}

func TestRmSpacesAtEnd( t *testing.T ){
  data := []struct{
    input    string
    output   string
  } {
    { "", "" },
    { " ", "" },
    { "\n\n\n", "" },
    { " \n\ta", " \n\ta" },
    { "line", "line" },
    { "line\n", "line" },
    { "line\t\v\rline", "line\t\v\rline" },
    { "line\t\v\nline\n", "line\t\v\nline" },
    { "  \nline\t\v\n\nline\n\n", "  \nline\t\v\n\nline" },
    { "line\t\v\n\n\tline\n\n\n", "line\t\v\n\n\tline" },
  }

  for _, d := range data {
    output := RmSpacesAtEnd( d.input )
    if output != d.output {
      t.Errorf( "RmSpacesAtEnd( %q ) \nreturn   %q\nexpected %q", d.input, output, d.output )
    }
  }
}

func TestRmSpacesAtStartup( t *testing.T ){
  data := []struct{
    input    string
    output   string
  } {
    { "", "" },
    { " ", "" },
    { "\n\n\n", "" },
    { " \n\ta", "a" },
    { "line", "line" },
    { "line\n", "line\n" },
    { "line\t\v\rline", "line\t\v\rline" },
    { "line\t\v\nline\n", "line\t\v\nline\n" },
    { "  \nline\t\v\n\nline\n\n", "line\t\v\n\nline\n\n" },
    { "line\t\v\n\n\tline\n\n\n", "line\t\v\n\n\tline\n\n\n" },
  }

  for _, d := range data {
    output := RmSpacesAtStartup( d.input )
    if output != d.output {
      t.Errorf( "RmSpacesAtStartup( %q ) \nreturn   %q\nexpected %q", d.input, output, d.output )
    }
  }
}

func TestRmSpacesToTheSides( t *testing.T ){
  data := []struct{
    input    string
    output   string
  } {
    { "", "" },
    { " ", "" },
    { "\n\n\n", "" },
    { " \n\ta", "a" },
    { "line", "line" },
    { "line\n", "line" },
    { "line\t\v\rline", "line\t\v\rline" },
    { "line\t\v\nline\n", "line\t\v\nline" },
    { "  \nline\t\v\n\nline\n\n", "line\t\v\n\nline" },
    { "line\t\v\n\n\tline\n\n\n", "line\t\v\n\n\tline" },
  }

  for _, d := range data {
    output := RmSpacesToTheSides( d.input )
    if output != d.output {
      t.Errorf( "RmSpacesToTheSides( %q ) \nreturn   %q\nexpected %q", d.input, output, d.output )
    }
  }
}

func TestLinelize( t *testing.T ){
  data := []struct{
    input    string
    output   string
  } {
    { "", "" },
    { " ", "" },
    { " hola\t\n", "hola" },
    { "\n\n\n", "" },
    { " \n\ta", "a" },
    { " \t\ta", "a" },
    { " \t\ta\ta", "a\ta" },
    { "line", "line" },
    { "line\n", "line" },
    { "line\t\v\rline", "line\t\v\rline" },
    { "line\t\v\nline\n++", "line line ++" },
    { "  \nline\t\v\n\nline\n\n", "line line" },
    { "line\t\v\n\n\tline\n\n\n", "line line" },
    { "  \nline\t\v\n \nline\n\n", "line line" },
  }

  for _, d := range data {
    output := Linelize( d.input )
    if output != d.output {
      t.Errorf( "Linelize( %q ) \nreturn   %q\nexpected %q", d.input, output, d.output )
    }
  }
}

func TestSpaceSwap( t *testing.T ){
  data := []struct{
    input    string
    output   string
  } {
    { "", "" },
    { " ", " " },
    { "\n\n\n", " " },
    { " \n\ta", " a" },
    { " \t\ta", " a" },
    { " \t\ta\ta", " a a" },
    { "line", "line" },
    { "line\n", "line " },
    { " hola\t\n", " hola " },
    { " hola  \t\n"," hola " },
    { "line\t\v\rline", "line line" },
    { "line\t\v\nline\n", "line line " },
    { "  \nline\t\v\n\nline\n\n", " line line " },
    { "line\t\v\n\n\tline\n\n\n", "line line " },
  }

  for _, d := range data {
    output := SpaceSwap( d.input, " " )
    if output != d.output {
      t.Errorf( "SpaceSwap( %q, %q ) \nreturn   %q\nexpected %q", d.input, " ", output, d.output )
    }
  }
}

func TestCountSpacesRegions( t *testing.T ){
  data := []struct{
    input  string
    output int
  } {
    { "", 0 },
    { "text", 0 },
    { " ", 1 },
    { "\n\n\n", 1 },
    { " \n\ta", 1 },
    { " \t\ta", 1 },
    { " \t\ta\ta", 2 },
    { "line\n", 1 },
    { "line\t\v\rline", 1 },
    { "line\t\v\nline\n", 2 },
    { "  \nline\t\v\n\nline\n\n", 3 },
    { "line\t\v\n\n\tline\n\n\n", 2 },
  }

  for _, d := range data {
    output := countSpacesRegions( d.input )
    if output != d.output {
      t.Errorf( "countSpacesRegions( %q ) \nreturn   %d\nexpected %d", d.input, output, d.output )
    }
  }
}


func TestRmIndent( t *testing.T ){
  data := []struct{
    input    string
    output   string
    n        int
  } {
    { "", "", 0 },
    { "", "", 5 },
    { "", "", 5 },
    { "  hola\n  hey", "  hola\n  hey", 5 },
    { "  hola\n  hey", "hola\nhey", 2 },
    { "   hola\n   hey\n", " hola\n hey\n", 2 },
    { "   hola\n   hey\n   hoy\nhi", " hola\n hey\n hoy\nhi", 2 },
    { "   hola\n   hey\n   hoy\nhi", "hola\nhey\nhoy\nhi", 3 },
  }

  for _, d := range data {
    output := RmIndent( d.input, d.n )
    if output != d.output {
      t.Errorf( "RmIndent( %q, %d ) \nreturn   %q\nexpected %q", d.n, d.input, output, d.output )
    }
  }
}

func TestCountIndentSpaces( t *testing.T ){
  data := []struct{
    input    string
    output   int
  } {
    { "", 0 },
    { "a", 0 },
    { "hola", 0 },
    { " ", 1 },
    { "\t\t ", 3 },
    { "\t \t ", 4 },
    { "\t \t hola", 4 },
    { "hola\t \t hola", 0 },
  }

  for _, d := range data {
    output := CountIndentSpaces( d.input )
    if output != d.output {
      t.Errorf( "CountIndentSpaces( %q ) \nreturn   %d\nexpected %d", d.input, output, d.output )
    }
  }
}

func TestDragTextByIndent( t *testing.T ){
  data := []struct{
    input    string
    n        int
    output   string
    l        int
  } {
    { "", 0, "", 0 },
    { "", 4, "", 0 },
    { "hola", 4, "", 0 },
    { "hola", 0, "hola", 4 },
    { "hola\nhi\n hoy", 0, "hola\nhi\n hoy", 12 },
    { "  hola\n   hi\n hoy", 2, "  hola\n   hi\n", 13 },
    { "  hola\n   hi\n\n   hoy", 2, "  hola\n   hi\n", 13 },
    { `  Nullamの私はあなたの意見を聞いた。 Fusce suscipit、wisi nec facilisis
  facilisis、est dui fermentum leo、欲求不満の欲求が欲しかった。 Nunc portaはテ
  ルトゥスです。 Nunc rutrum turpis sed pede。 bibendumになる Aliquam posuere。
  Nunc aliquet、augue nec adipiscing interdum、lacus tellus malesuada massa、お
  よび他のものは、他のものとは異なる。 Pellentesque condimentum、magna ut
  susipipit hendrerit、ipsum augue ornare nulla、non luctus diam neque sit amet
  urna。 Curabiturは前庭ロレムを訴える。 虚脱、リベロの非暴行、巨大orciの痛み、
  nepta nea naclacinia eros。 Sid id ligulaはest convallis temporです。
  Curabitur lacinia pulvinar nibh。 サピエンナム。
counte
out
`, 2, `  Nullamの私はあなたの意見を聞いた。 Fusce suscipit、wisi nec facilisis
  facilisis、est dui fermentum leo、欲求不満の欲求が欲しかった。 Nunc portaはテ
  ルトゥスです。 Nunc rutrum turpis sed pede。 bibendumになる Aliquam posuere。
  Nunc aliquet、augue nec adipiscing interdum、lacus tellus malesuada massa、お
  よび他のものは、他のものとは異なる。 Pellentesque condimentum、magna ut
  susipipit hendrerit、ipsum augue ornare nulla、non luctus diam neque sit amet
  urna。 Curabiturは前庭ロレムを訴える。 虚脱、リベロの非暴行、巨大orciの痛み、
  nepta nea naclacinia eros。 Sid id ligulaはest convallis temporです。
  Curabitur lacinia pulvinar nibh。 サピエンナム。
`, 781 },
  }

  for _, d := range data {
    output, l := DragTextByIndent( d.input, d.n )
    if output != d.output || l != d.l {
      t.Errorf( "DragTextByIndent( %q, %d ) \nreturn   [%d] %q\nexpected [%d] %q", d.input, d.n, l, output, d.l, d.output )
    }
  }
}

func TestDragLineAndTextByIndent( t *testing.T ){
  data := []struct{
    input    string
    n        int
    output   string
    l        int
  } {
    { "", 0, "", 0 },
    { "", 4, "", 0 },
    { "hola", 4, "hola", 4 },
    { "hola", 0, "hola", 4 },
    { "hola\nhi\n hoy", 0, "hola\nhi\n hoy", 12 },
    { "hola\n  hi\n  hoy", 2, "hola\n  hi\n  hoy", 15 },
    { "hola\nhi\n  hoy", 2, "hola\n", 5 },
    { "hola\nhi\nhoy", 2, "hola\n", 5 },
    { "hola\n  ni\nhoy", 2, "hola\n  ni\n", 10 },
    { "  hola\n   hi\n hoy", 2, "  hola\n   hi\n", 13 },
    { "  hola\n   hi\n\n   hoy", 2, "  hola\n   hi\n", 13 },
    { `  Nullamの私はあなたの意見を聞いた。 Fusce suscipit、wisi nec facilisis
  facilisis、est dui fermentum leo、欲求不満の欲求が欲しかった。 Nunc portaはテ
  ルトゥスです。 Nunc rutrum turpis sed pede。 bibendumになる Aliquam posuere。
  Nunc aliquet、augue nec adipiscing interdum、lacus tellus malesuada massa、お
  よび他のものは、他のものとは異なる。 Pellentesque condimentum、magna ut
  susipipit hendrerit、ipsum augue ornare nulla、non luctus diam neque sit amet
  urna。 Curabiturは前庭ロレムを訴える。 虚脱、リベロの非暴行、巨大orciの痛み、
  nepta nea naclacinia eros。 Sid id ligulaはest convallis temporです。
  Curabitur lacinia pulvinar nibh。 サピエンナム。
counte
out
`, 2, `  Nullamの私はあなたの意見を聞いた。 Fusce suscipit、wisi nec facilisis
  facilisis、est dui fermentum leo、欲求不満の欲求が欲しかった。 Nunc portaはテ
  ルトゥスです。 Nunc rutrum turpis sed pede。 bibendumになる Aliquam posuere。
  Nunc aliquet、augue nec adipiscing interdum、lacus tellus malesuada massa、お
  よび他のものは、他のものとは異なる。 Pellentesque condimentum、magna ut
  susipipit hendrerit、ipsum augue ornare nulla、non luctus diam neque sit amet
  urna。 Curabiturは前庭ロレムを訴える。 虚脱、リベロの非暴行、巨大orciの痛み、
  nepta nea naclacinia eros。 Sid id ligulaはest convallis temporです。
  Curabitur lacinia pulvinar nibh。 サピエンナム。
`, 781 },
  }

  for _, d := range data {
    output, l := DragLineAndTextByIndent( d.input, d.n )
    if output != d.output || l != d.l {
      t.Errorf( "DragTextByIndent( %q, %d ) \nreturn   [%d] %q\nexpected [%d] %q", d.input, d.n, l, output, d.l, d.output )
    }
  }
}

func TestDragAllTextByIndent( t *testing.T ){
  data := []struct{
    input    string
    n        int
    output   string
    l        int
  } {
    { "", 0, "", 0 },
    { "", 4, "", 0 },
    { "\n hola", 4, "\n", 1 },
    { "hola", 4, "", 0 },
    { "hola", 0, "hola", 4 },
    { "hola\nhi\n hoy", 0, "hola\nhi\n hoy", 12 },
    { "  hola\n   hi\n hoy", 2, "  hola\n   hi\n", 13 },
    { "  hola\n   hi\n\n   hoy", 2, "  hola\n   hi\n\n   hoy", 20 },
    { "  hola\n\n\n\t  hi\n\n   hoy", 2, "  hola\n\n\n\t  hi\n\n   hoy", 22 },
    { "\n\nhola\n\n\n\t  hi\n\n   hoy", 2, "\n\n", 2 },
    { "\n\n  hola\n\n\n\t  hi\n\n   hoy", 2, "\n\n  hola\n\n\n\t  hi\n\n   hoy", 24 },
    { "\n\n  hola\n\n\n\t  hi\n..hoy", 2, "\n\n  hola\n\n\n\t  hi\n", 17 },
    { `  01 Lorem ipsum es el texto que se usa habitualmente en diseño gráfico en
  demostraciones de tipografías o de borradores de diseño para probar el diseño
  visual antes de insertar el texto final.

..figure > 02`, 2,
`  01 Lorem ipsum es el texto que se usa habitualmente en diseño gráfico en
  demostraciones de tipografías o de borradores de diseño para probar el diseño
  visual antes de insertar el texto final.

`, 204 },

  }

  for _, d := range data {
    output, l := DragAllTextByIndent( d.input, d.n )
    if output != d.output || l != d.l {
      t.Errorf( "DragAllTextByIndent( %q, %d ) \nreturn   [%d] %q\nexpected [%d] %q", d.input, d.n, l, output, d.l, d.output )
    }
  }
}

//////// benchmarks

const ssIn  = "  \nline-a\t\v\n\nline-b\n\nline-c\nline-d\t\v\n\nline-en\n"
const ssSwp = "––"
const ssOut = "––line-a––line-b––line-c––line-d––line-en––"

var recom   = regexp.MustCompile( "[[:space:]]+" )

func SpaceSwapRegexp( str, swap string ) string {
  return recom.ReplaceAllString( str, swap )
}

func BenchmarkSpaceSwapRegexp( b *testing.B ){
  for i := 0; i < b.N; i++ {
    if( SpaceSwapRegexp( ssIn, ssSwp ) != ssOut ){
      b.Fatalf( "BenchmarkSpaceSwapRegexp(): no match" )
    }
  }
}


var reSpace = regexp4.Compile( "<:s+>" )

func SpaceSwapRegexp4( str, swap string ) string {
  res := reSpace.Copy()
  if res.FindString( str ) {
    return res.RplCatch( swap, 1 )
  }

  return str
}

func BenchmarkSpaceSwapRegexp4( b *testing.B ){
  for i := 0; i < b.N; i++ {
    if( SpaceSwapRegexp4( ssIn, ssSwp ) != ssOut ){
      b.Fatal( "BenchmarkSpaceSwapRegexp(): no match" )
    }
  }
}

func SpaceSwapBuffer( str, swap string ) string {
  var k bytes.Buffer

  for i := 0; i < len( str );  {
    switch str[i] {
    case ' ', '\t', '\n', '\v', '\f', '\r' :
      i += CountInitSpaces( str[i:] )
      k.WriteString( swap )
    default: k.WriteByte( str[i] ); i++
    }
  }

  return k.String()
}

func BenchmarkSpaceSwapBuffer( b *testing.B ){
  for i := 0; i < b.N; i++ {
    if( SpaceSwapBuffer( ssIn, ssSwp ) != ssOut ){
      b.Fatal( "BenchmarkSpaceSwapRegexp(): no match" )
    }
  }
}

func SpaceSwapString( str, swap string ) string {
  k := ""

  for i := 0; i < len( str );  {
    switch str[i] {
    case ' ', '\t', '\n', '\v', '\f', '\r' :
      i += CountInitSpaces( str[i:] )
      k += swap
    default: k += str[i:i+1]; i++
    }
  }

  return k
}

func BenchmarkSpaceSwapString( b *testing.B ){
  for i := 0; i < b.N; i++ {
    if( SpaceSwapString( ssIn, ssSwp ) != ssOut ){
      b.Fatal( "BenchmarkSpaceSwapRegexp(): no match" )
    }
  }
}

func SpaceSwapFor( str, swap string ) string {
  i, j, k := 0, 0, make( []byte, len( str ) + len( swap ) * countSpacesRegions( str ) )

  for ; i < len( str );  {
    switch str[i] {
    case ' ', '\t', '\n', '\v', '\f', '\r' :
      i += CountInitSpaces( str[i:] )
      for h := 0; h < len( swap ); h++ {
        k[j] = swap[h]; j++
      }
    default: k[j] = str[i]; j++; i++
    }
  }

  return string( k[:j] )
}

func BenchmarkSpaceSwapFor( b *testing.B ){
  for i := 0; i < b.N; i++ {
    if( SpaceSwapFor( ssIn, ssSwp ) != ssOut ){
      b.Fatal( "BenchmarkSpaceSwapRegexp(): no match" )
    }
  }
}

func BenchmarkSpaceSwap( b *testing.B ){
  for i := 0; i < b.N; i++ {
    if( SpaceSwap( ssIn, ssSwp ) != ssOut ){
      b.Fatal( "BenchmarkSpaceSwapRegexp(): no match" )
    }
  }
}
