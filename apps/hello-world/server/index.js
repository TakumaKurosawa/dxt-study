import { Server } from '@modelcontextprotocol/sdk/server/index.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import {
  ListToolsRequestSchema,
  CallToolRequestSchema,
} from '@modelcontextprotocol/sdk/types.js';

// MCPサーバーの作成
const server = new Server(
  {
    name: 'Hello World Server',
    version: '1.0.0',
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

// 利用可能なツール一覧を返すハンドラー
server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: 'say_hello',
        description: '指定された名前に対して挨拶を返します',
        inputSchema: {
          type: 'object',
          properties: {
            name: {
              type: 'string',
              description: '挨拶する相手の名前',
            },
            language: {
              type: 'string',
              description: '挨拶の言語 (japanese, english)',
              enum: ['japanese', 'english'],
              default: 'japanese',
            },
          },
          required: ['name'],
        },
      },
      {
        name: 'get_time',
        description: '現在の時刻を返します',
        inputSchema: {
          type: 'object',
          properties: {
            format: {
              type: 'string',
              description: '時刻のフォーマット (12h, 24h)',
              enum: ['12h', '24h'],
              default: '24h',
            },
          },
        },
      },
    ],
  };
});

// ツール実行のハンドラー
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;

  switch (name) {
    case 'say_hello': {
      const language = args.language || 'japanese';
      const greeting =
        language === 'english'
          ? `Hello, ${args.name}!`
          : `こんにちは、${args.name}さん！`;

      return {
        content: [
          {
            type: 'text',
            text: greeting,
          },
        ],
      };
    }

    case 'get_time': {
      const now = new Date();
      const format = args.format || '24h';

      let timeString;
      if (format === '12h') {
        timeString = now.toLocaleTimeString('ja-JP', {
          hour12: true,
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
        });
      } else {
        timeString = now.toLocaleTimeString('ja-JP', {
          hour12: false,
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
        });
      }

      return {
        content: [
          {
            type: 'text',
            text: `現在の時刻: ${timeString}`,
          },
        ],
      };
    }

    default:
      throw new Error(`Unknown tool: ${name}`);
  }
});

// エラーハンドリング
server.onerror = (error) => {
  console.error('[MCP Error]', error);
};

// サーバーの起動
async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.log('Hello World MCP server running on stdio');
}

main().catch((error) => {
  console.error('Failed to start server:', error);
  process.exit(1);
});
