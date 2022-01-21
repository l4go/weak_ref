# type WeakRef
弱参照の機能を提供するモジュールです。  
atomic操作で排他制御を行っているので、Lock-freeです。

## import

```go
import "github.com/l4go/weak_ref"
```
vendoringして使うことを推奨します。

## 利用サンプル

* [弱参照のサンプル](../examples/ex_weak_ref/ex_weak_ref.go)
* [所有権譲渡のサンプル](../examples/ex_ref_move/ex_ref_move.go)
* [CAS処理のサンプル](../examples/ex_ref_cas/ex_ref_cas.go)
* [CAS処理のループを自前で書いたサンプル](../examples/ex_ref_cas2/ex_ref_cas2.go)

## 参照の扱い
\*WeakRefでの参照の扱いは以下の通りです。

* 参照は、ポインタ型の値で表現する。
* 内部へ保存された参照へのアクセスは、atomic操作で行う。
* 参照を無効化でき、nil値を無効状態を示すものとして扱う。
* 入力では、参照のポインタ型の値を`interface{}`で受け取る。
* 出力では、参照のポインタ型の値を`interface{}`として返す。
* 初期化時(New()時)に、参照の型を固定する。(静的な型アサーションが利用可能)

## メソッド概略

### func New(val interface{}) \*WeakRef
\*WeakRefを生成します。
参照の利用終了時に、無効にする(Reset()等)ことで弱参照として動作します。  
初期値を使って、値の設定だけでなく、参照の型を固定化も行います。初期化後は、初期値と異なる型の値を入力するとpanicします。

初期値だけは型情報が必要なので、nil型のnil値が利用できません。無効状態(nil値)に初期化したい場合は、以下の例のように型ありのnil値を利用してください。

ポインタ型のnil値(無効状態)で初期化する例

```go 
type Test struct {
	val int
}

var wr =  weak_ref.New((*Test)(nil))
```

### func (wr \*WeakRef) Move() \*WeakRef
参照の所有権を譲渡した\*WeakRefを生成します。  
元となった\*WeakRefは無効にされます。

別の処理(関数など)へ、参照の管理ごと譲渡する用途を想定しています。

### func (wr \*WeakRef) Reset()
参照を無効にします。  
弱参照として動作させるには、参照の利用が終了した際に、Reset()の呼び出しが必要です。

### func (wr \*WeakRef) Get() interface{}
参照の値を取得します。

### func (wr \*WeakRef) Set(val interface{})
参照の値を変更します。   
静的なnilで無効化する場合は、簡素に記述できるReset()を利用するべきです。  
goroutine間での排他制御が必要な場合は、CasUpdate()を利用してください。

### func (wr \*WeakRef) Swap(new\_val interface{}) interface{}
参照の値を更新し、変更前の値を返します。  
goroutine間での排他制御が必要な場合は、CasUpdate()を利用してください。

### func (wr \*WeakRef) CompareAndSwap(old\_val, new\_val interface{}) bool
CAS(コンペア・アンド・スワップ、Compare-and-swap)での、参照の値の更新を試み、
成功した場合はtureを、失敗した場合はfalseを返します。

通常の用途では、より簡素に記述できるCasUpdate()を利用すべきです。

### func (wr \*WeakRef) CasUpdate(f CasUpdateFunc) interface{}
CASアルゴリズム(CAS方式のループを実施)で、参照の値を更新します。  
変更後の値の計算方法をCasUpdateFunc型の関数で渡します。  
参照が無効な場合は失敗し、nil値を返します。
更新が成功した場合は変更後の値を返します。

### type CasUpdateFunc func(v interface{}) interface{}
CasUpdate()で利用する、古い値から新しい値への計算方法を記述する関数の型です。  
CASアルゴリズムではコリジョン時に再計算が必要なため、CasUpdateFunc型の関数は、1回の更新処理で複数回呼ばれることがあります。

インクリメントするCasUpdateFunc型の関数の例(`func incr()`)

```go
type Test struct {
        V int
}

func incr(v interface{}) interface{} {
        return &Test{V: v.V+1}
}
```
