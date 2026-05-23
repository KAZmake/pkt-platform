'use client';

import { useState, useRef, useEffect, useCallback } from 'react';
import Link from 'next/link';

interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string;
}

const WELCOME: Message = {
  id: 'welcome',
  role: 'assistant',
  content:
    'Здравствуйте! Я виртуальный ассистент ПКТ. Готов ответить на вопросы о программах кредитования, условиях займов и процедуре подачи заявки.',
};

let msgCounter = 0;
function newId() {
  return `msg-${++msgCounter}`;
}

export function ChatWidget() {
  const [open, setOpen] = useState(false);
  const [messages, setMessages] = useState<Message[]>([WELCOME]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [hasNew, setHasNew] = useState(false);

  const bottomRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // Auto-scroll on new messages
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, loading]);

  // Focus textarea on open, clear badge
  useEffect(() => {
    if (open) {
      textareaRef.current?.focus();
      setHasNew(false);
    }
  }, [open]);

  const userCount = messages.filter((m) => m.role === 'user').length;

  const send = useCallback(async () => {
    const text = input.trim();
    if (!text || loading) return;

    const userMsg: Message = { id: newId(), role: 'user', content: text };
    setMessages((prev) => [...prev, userMsg]);
    setInput('');
    // Reset textarea height
    if (textareaRef.current) textareaRef.current.style.height = 'auto';
    setLoading(true);

    try {
      const history = [...messages, userMsg].map((m) => ({
        role: m.role,
        content: m.content,
      }));
      const res = await fetch('/api/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ messages: history }),
      });

      const data: { content?: string; error?: string } = await res.json();

      setMessages((prev) => [
        ...prev,
        {
          id: newId(),
          role: 'assistant',
          content: data.content ?? 'Извините, не удалось получить ответ. Попробуйте позже.',
        },
      ]);
      if (!open) setHasNew(true);
    } catch {
      setMessages((prev) => [
        ...prev,
        {
          id: newId(),
          role: 'assistant',
          content: 'Произошла ошибка связи. Попробуйте позже или свяжитесь с нами по телефону.',
        },
      ]);
    } finally {
      setLoading(false);
    }
  }, [input, loading, messages, open]);

  function handleKeyDown(e: React.KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      send();
    }
  }

  function handleInput(e: React.FormEvent<HTMLTextAreaElement>) {
    const el = e.currentTarget;
    el.style.height = 'auto';
    el.style.height = `${Math.min(el.scrollHeight, 96)}px`;
  }

  return (
    <>
      {/* Chat panel */}
      {open && (
        <div
          className="fixed bottom-20 right-4 z-50 w-80 sm:w-96 bg-white rounded-2xl shadow-2xl border border-gray-200 flex flex-col overflow-hidden"
          style={{ height: '500px' }}
        >
          {/* Header */}
          <div className="flex items-center gap-3 px-4 py-3 bg-brand-green shrink-0">
            <div className="w-8 h-8 bg-white/20 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0">
              ИИ
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-white font-semibold text-sm">Ассистент ПКТ</p>
              <div className="flex items-center gap-1.5">
                <span className="w-1.5 h-1.5 bg-emerald-400 rounded-full" />
                <p className="text-white/70 text-xs">Онлайн</p>
              </div>
            </div>
            <button
              onClick={() => setOpen(false)}
              className="text-white/60 hover:text-white transition-colors text-lg leading-none shrink-0"
              aria-label="Закрыть"
            >
              ✕
            </button>
          </div>

          {/* Messages */}
          <div className="flex-1 overflow-y-auto px-4 py-3 space-y-3">
            {messages.map((msg) => (
              <div
                key={msg.id}
                className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
              >
                <div
                  className={`max-w-[82%] rounded-2xl px-3.5 py-2.5 text-sm leading-relaxed whitespace-pre-wrap break-words ${
                    msg.role === 'user'
                      ? 'bg-brand-green text-white rounded-tr-sm'
                      : 'bg-gray-100 text-gray-800 rounded-tl-sm'
                  }`}
                >
                  {msg.content}
                </div>
              </div>
            ))}

            {/* Typing indicator */}
            {loading && (
              <div className="flex justify-start">
                <div className="bg-gray-100 rounded-2xl rounded-tl-sm px-4 py-3">
                  <div className="flex gap-1 items-center h-4">
                    {[0, 0.15, 0.3].map((delay, i) => (
                      <span
                        key={i}
                        className="w-1.5 h-1.5 bg-gray-400 rounded-full animate-bounce"
                        style={{ animationDelay: `${delay}s` }}
                      />
                    ))}
                  </div>
                </div>
              </div>
            )}
            <div ref={bottomRef} />
          </div>

          {/* Create ticket CTA — shown after 3+ user messages */}
          {userCount >= 3 && (
            <div className="px-4 py-2 border-t border-gray-100 shrink-0">
              <Link href="/cabinet/tickets" className="text-xs text-brand-green hover:underline">
                Не нашли ответ? Создать обращение →
              </Link>
            </div>
          )}

          {/* Input */}
          <div className="px-3 pb-3 pt-2 border-t border-gray-100 shrink-0">
            <div className="flex items-end gap-2">
              <textarea
                ref={textareaRef}
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={handleKeyDown}
                onInput={handleInput}
                rows={1}
                maxLength={1000}
                placeholder="Напишите вопрос..."
                disabled={loading}
                className="flex-1 rounded-xl border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-green resize-none overflow-y-auto"
                style={{ minHeight: '36px', maxHeight: '96px' }}
              />
              <button
                onClick={send}
                disabled={loading || !input.trim()}
                aria-label="Отправить"
                className="w-9 h-9 bg-brand-green text-white rounded-xl flex items-center justify-center hover:bg-brand-green-dark disabled:opacity-40 shrink-0 transition-colors"
              >
                <svg
                  width="16"
                  height="16"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  strokeWidth={2}
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"
                  />
                </svg>
              </button>
            </div>
            <p className="text-xs text-gray-300 mt-1.5 text-right">
              Enter — отправить · Shift+Enter — новая строка
            </p>
          </div>
        </div>
      )}

      {/* Floating button */}
      <button
        onClick={() => setOpen((v) => !v)}
        aria-label={open ? 'Закрыть ассистента' : 'Открыть ассистента'}
        className="fixed bottom-4 right-4 z-50 w-14 h-14 bg-brand-green text-white rounded-full shadow-lg hover:bg-brand-green-dark transition-colors flex items-center justify-center"
      >
        {open ? (
          <svg
            width="20"
            height="20"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth={2.5}
          >
            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        ) : (
          <svg
            width="24"
            height="24"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth={2}
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"
            />
          </svg>
        )}
        {hasNew && !open && (
          <span className="absolute top-0.5 right-0.5 w-4 h-4 bg-red-500 rounded-full text-xs flex items-center justify-center text-white font-bold leading-none">
            !
          </span>
        )}
      </button>
    </>
  );
}
