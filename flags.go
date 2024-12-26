// Copyright 2024 tobyzxj
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package coap

// (uint32)Flags:         标志位
//                                               3                       2                       1                       0
//                         31 30 29 28 27 26 25 24 23 22 21 20 19 18 17 16 15 14 13 12 11 10  9  8  7  6  5  4  3  2  1  0
//                        +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//                        | X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| X| 0|
//                        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//                         bit0  : 是否需要平台回复，0: 需要回复，1: 不需要回复
//                         others: 保留

const (
	FlagNoAck = 0x00000001
)

// FlagIsNoAck: 判断是否不需要回复
func FlagIsNoAck(flags uint32) bool {
	return flags&FlagNoAck != 0
}
