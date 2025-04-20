import { defineConfig } from '@kubb/core'
import { pluginOas } from '@kubb/plugin-oas'
import { pluginTs } from '@kubb/plugin-ts'

export default defineConfig(() => {
  return {
    root: '.',
    input: {
      path: '../backend/docs/swagger.json',
    },
    output: {
      path: './src/types',
    },
    plugins: [
      pluginOas(),
      pluginTs({
        output: {
          path: './gen',
        },
        exclude: [
          {
            type: 'tag',
            pattern: 'store',
          },
        ],
        enumType: 'asConst',
        enumSuffix: 'Enum',
        dateType: 'date',
        unknownType: 'unknown',
        oasType: 'infer',
      }),
    ],
  }
})
