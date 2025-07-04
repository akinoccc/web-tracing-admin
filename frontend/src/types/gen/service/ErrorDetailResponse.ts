/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { ServiceErrorEventItem } from './ErrorEventItem.ts'
import type { ServiceErrorGroupItem } from './ErrorGroupItem.ts'

export type ServiceErrorDetailResponse = {
  /**
   * @type array | undefined
   */
  events?: ServiceErrorEventItem[]
  /**
   * @type object | undefined
   */
  group?: ServiceErrorGroupItem
  /**
   * @type integer | undefined
   */
  total?: number
}