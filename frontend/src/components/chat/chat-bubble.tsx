import { useMemo } from "react";
import ReactMarkdown from "react-markdown";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { oneDark } from "react-syntax-highlighter/dist/esm/styles/prism";
import type { Message } from "@/lib/store/chat";
import { cn } from "@/lib/utils";

interface ChatBubbleProps {
  message: Message;
  isStreaming?: boolean;
}

export function ChatBubble({ message, isStreaming }: ChatBubbleProps) {
  const isUser = message.role === "user";

  const components = useMemo(
    () => ({
      code({
        className,
        children,
        ...props
      }: {
        className?: string;
        children?: React.ReactNode;
      }) {
        const match = /language-(\w+)/.exec(className ?? "");
        const code = String(children).replace(/\n$/, "");
        if (match) {
          return (
            <SyntaxHighlighter
              style={oneDark}
              language={match[1]}
              PreTag="div"
              customStyle={{
                margin: 0,
                borderRadius: "0.375rem",
                fontSize: "0.8125rem",
                fontFamily: '"JetBrains Mono", monospace',
              }}
            >
              {code}
            </SyntaxHighlighter>
          );
        }
        return (
          <code
            className="rounded bg-code-block px-1 py-0.5 font-mono text-sm text-text-primary"
            {...props}
          >
            {children}
          </code>
        );
      },
    }),
    []
  );

  return (
    <div
      className={cn(
        "flex w-full gap-3 px-4 py-3",
        isUser ? "justify-end" : "justify-start"
      )}
    >
      {!isUser && (
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-accent/20 font-mono text-xs text-accent">
          IG
        </div>
      )}
      <div
        className={cn(
          "max-w-[75%] rounded-xl px-4 py-3 text-sm leading-relaxed",
          isUser
            ? "bg-user-bubble text-white"
            : "bg-surface border border-border text-text-primary"
        )}
      >
        {isUser ? (
          <p className="whitespace-pre-wrap">{message.content}</p>
        ) : (
          <div className="prose prose-sm prose-invert max-w-none [&_pre]:bg-transparent [&_pre]:p-0">
            <ReactMarkdown components={components}>
              {message.content}
            </ReactMarkdown>
            {isStreaming && (
              <span className="ml-0.5 inline-block h-4 w-1 animate-pulse bg-accent" />
            )}
          </div>
        )}
      </div>
      {isUser && (
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-accent font-mono text-xs text-white">
          U
        </div>
      )}
    </div>
  );
}
