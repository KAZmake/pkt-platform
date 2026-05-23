'use client';

interface GrafanaFrameProps {
  src: string;
}

export function GrafanaFrame({ src }: GrafanaFrameProps) {
  return (
    <iframe
      src={src}
      className="w-full h-full rounded-xl border border-gray-200"
      title="Grafana Analytics"
      sandbox="allow-scripts allow-same-origin allow-forms allow-popups"
      referrerPolicy="no-referrer-when-downgrade"
    />
  );
}
