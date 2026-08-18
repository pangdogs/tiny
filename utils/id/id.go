/*
 * This file is part of Golaxy Distributed Service Development Framework.
 *
 * Golaxy Distributed Service Development Framework is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * Golaxy Distributed Service Development Framework is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with Golaxy Distributed Service Development Framework. If not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (c) 2024 pangdogs.
 */

package id

import "strconv"

// Nil 表示未分配的 ID。
const Nil ID = 0

// ID 表示 Runtime 内的本地整数 ID。
type ID int64

// IsNil 报告 ID 是否尚未分配。
func (id ID) IsNil() bool {
	return id == Nil
}

// String 返回十进制文本。
func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}
