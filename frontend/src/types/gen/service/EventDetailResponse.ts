/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { ModelClickEvent } from '../model/ClickEvent.ts'
import type { ModelErrorEvent } from '../model/ErrorEvent.ts'
import type { ModelEvent } from '../model/Event.ts'
import type { ModelExposureEvent } from '../model/ExposureEvent.ts'
import type { ModelPerformanceEvent } from '../model/PerformanceEvent.ts'
import type { ModelRequestEvent } from '../model/RequestEvent.ts'
import type { ModelRouteEvent } from '../model/RouteEvent.ts'

export type ServiceEventDetailResponse = {
  /**
   * @type object | undefined
   */
  click?: ModelClickEvent
  /**
   * @type object | undefined
   */
  errorEvent?: ModelErrorEvent
  /**
   * @type object | undefined
   */
  event?: ModelEvent
  /**
   * @type object | undefined
   */
  exposure?: ModelExposureEvent
  /**
   * @type object | undefined
   */
  performance?: ModelPerformanceEvent
  /**
   * @type object | undefined
   */
  request?: ModelRequestEvent
  /**
   * @type object | undefined
   */
  route?: ModelRouteEvent
}