#!/usr/bin/env node
// .github/scripts/generate-dashboard.js
// Парсит _docs/roadmap.md и генерирует docs-site/index.html

const fs = require('fs');
const path = require('path');

const roadmapPath = path.join(process.cwd(), '_docs', 'roadmap.md');
const outDir = path.join(process.cwd(), 'docs-site');
const outFile = path.join(outDir, 'index.html');

if (!fs.existsSync(outDir)) fs.mkdirSync(outDir, { recursive: true });

const raw = fs.readFileSync(roadmapPath, 'utf8');
const lines = raw.split('\n');

// ── Парсим блок ТЕКУЩИЙ СТАТУС ──────────────────────────────────────────────
const statusBlock = {};
const statusMatch = raw.match(/```\n([\s\S]*?)```/);
if (statusMatch) {
  statusMatch[1].split('\n').forEach(line => {
    const m = line.match(/^(.+?):\s+(.+)$/);
    if (m) statusBlock[m[1].trim()] = m[2].trim();
  });
}

// ── Парсим фазы и задачи ────────────────────────────────────────────────────
const phases = [];
let currentPhase = null;
let currentSubphase = null;

const phaseRe  = /^## (Фаза \d+)\s*[—–]\s*(.+?)\s*\(([^)]+)\)\s*`\[(.)\]/;
const subRe    = /^### (\d+\.\d+)\.\s+(.+)/;
const taskRe   = /^\|\s*(\d+[\.\d]*)\s*\|\s*(.+?)\s*\|\s*`\[(.)\]`\s*\|/;

lines.forEach(line => {
  const pm = line.match(phaseRe);
  if (pm) {
    currentPhase = {
      id: pm[1],
      title: pm[2].replace(/\s*`[^`]*`\s*/g, '').trim(),
      duration: pm[3],
      done: pm[4] === 'x',
      subphases: [],
      tasks: [],
    };
    currentSubphase = null;
    phases.push(currentPhase);
    return;
  }
  const sm = line.match(subRe);
  if (sm && currentPhase) {
    currentSubphase = { id: sm[1], title: sm[2], tasks: [] };
    currentPhase.subphases.push(currentSubphase);
    return;
  }
  const tm = line.match(taskRe);
  if (tm && currentPhase) {
    const task = { num: tm[1], title: tm[2].replace(/\*\*/g,'').split('|')[0].trim(), done: tm[3] === 'x' };
    if (currentSubphase) currentSubphase.tasks.push(task);
    else currentPhase.tasks.push(task);
  }
});

// ── Вычисляем прогресс ──────────────────────────────────────────────────────
function phaseStats(phase) {
  const allTasks = [
    ...phase.tasks,
    ...phase.subphases.flatMap(s => s.tasks),
  ];
  const total = allTasks.length;
  const done  = allTasks.filter(t => t.done).length;
  return { total, done, pct: total ? Math.round((done / total) * 100) : 0 };
}

const totalTasks = phases.reduce((acc, p) => acc + phaseStats(p).total, 0);
const doneTasks  = phases.reduce((acc, p) => acc + phaseStats(p).done, 0);
const overallPct = totalTasks ? Math.round((doneTasks / totalTasks) * 100) : 0;

const currentPhaseName = statusBlock['Фаза'] || '—';
const lastDone   = statusBlock['Последнее'] || '—';
const nextStep   = statusBlock['Следующий шаг'] || '—';
const blocked    = statusBlock['Заблокировано'] || '—';
const status     = statusBlock['Статус'] || '—';

// ── Цвета фаз ───────────────────────────────────────────────────────────────
const phaseColors = [
  '#1B4F8A','#166534','#1B4F8A','#4C1D95',
  '#134E4A','#92400E','#1B4F8A','#134E4A',
  '#4C1D95','#1B4F8A',
];
const phaseBg = [
  '#D6E4F5','#DCFCE7','#D6E4F5','#EDE9FE',
  '#CCFBF1','#FEF3C7','#D6E4F5','#CCFBF1',
  '#EDE9FE','#D6E4F5',
];

// ── Генерируем карточки фаз ─────────────────────────────────────────────────
function renderTask(t) {
  const icon = t.done
    ? `<span class="task-icon done">✓</span>`
    : `<span class="task-icon todo">○</span>`;
  return `<div class="task ${t.done ? 'task-done' : ''}">${icon}<span>${escHtml(t.title)}</span></div>`;
}

function renderSubphase(sp) {
  const done = sp.tasks.filter(t=>t.done).length;
  const total = sp.tasks.length;
  return `
    <div class="subphase">
      <div class="subphase-header">
        <span class="subphase-title">${escHtml(sp.title)}</span>
        <span class="subphase-count">${done}/${total}</span>
      </div>
      ${sp.tasks.map(renderTask).join('')}
    </div>`;
}

function renderPhase(p, i) {
  const { total, done, pct } = phaseStats(p);
  const color = phaseColors[i] || '#1B4F8A';
  const bg    = phaseBg[i]    || '#D6E4F5';
  const statusLabel = done === total && total > 0
    ? `<span class="badge badge-done">Готово</span>`
    : done > 0
      ? `<span class="badge badge-active">В процессе</span>`
      : `<span class="badge badge-todo">Не начато</span>`;

  const allFlat = p.subphases.length === 0;

  return `
  <div class="phase-card" style="--phase-color:${color};--phase-bg:${bg}">
    <div class="phase-header">
      <div class="phase-meta">
        <span class="phase-id">${p.id}</span>
        <div>
          <div class="phase-title">${escHtml(p.title)}</div>
          <div class="phase-duration">${escHtml(p.duration)}</div>
        </div>
      </div>
      <div class="phase-right">
        ${statusLabel}
        <span class="phase-pct">${pct}%</span>
      </div>
    </div>
    <div class="phase-progress-bar"><div class="phase-progress-fill" style="width:${pct}%"></div></div>
    <div class="phase-body">
      ${allFlat
        ? p.tasks.map(renderTask).join('')
        : p.subphases.map(renderSubphase).join('') + p.tasks.map(renderTask).join('')
      }
    </div>
    <div class="phase-footer">${done} из ${total} задач выполнено</div>
  </div>`;
}

function escHtml(s) {
  return String(s)
    .replace(/&/g,'&amp;')
    .replace(/</g,'&lt;')
    .replace(/>/g,'&gt;')
    .replace(/"/g,'&quot;');
}

const now = new Date().toLocaleString('ru-RU', { timeZone: 'Asia/Oral', dateStyle: 'long', timeStyle: 'short' });

// ── Собираем HTML ────────────────────────────────────────────────────────────
const html = `<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>ТОО «Первое кредитное товарищество» — Roadmap</title>
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

  :root {
    --blue: #1B4F8A;
    --blue-light: #D6E4F5;
    --dark: #1A1A2E;
    --gray: #F4F6FA;
    --border: #E2EAF4;
    --text-secondary: #6B7A99;
    --green: #166534;
    --green-bg: #DCFCE7;
    --amber: #92400E;
    --amber-bg: #FEF3C7;
    --radius: 12px;
  }

  body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
    background: #EEF2F9;
    color: var(--dark);
    min-height: 100vh;
  }

  /* ── HEADER ── */
  .site-header {
    background: var(--blue);
    padding: 24px 32px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 16px;
  }
  .site-title { color: #fff; font-size: 20px; font-weight: 700; }
  .site-sub   { color: rgba(255,255,255,0.7); font-size: 13px; margin-top: 2px; }
  .updated    { color: rgba(255,255,255,0.6); font-size: 12px; }

  /* ── OVERALL PROGRESS ── */
  .overall {
    background: #fff;
    border-bottom: 1px solid var(--border);
    padding: 20px 32px;
    display: grid;
    grid-template-columns: 1fr auto;
    gap: 24px;
    align-items: center;
  }
  .overall-bar-wrap { width: 100%; }
  .overall-label { font-size: 13px; color: var(--text-secondary); margin-bottom: 8px; }
  .overall-bar {
    height: 12px; border-radius: 99px;
    background: var(--blue-light); overflow: hidden;
  }
  .overall-fill {
    height: 100%; border-radius: 99px;
    background: var(--blue);
    transition: width 0.6s ease;
  }
  .overall-pct {
    font-size: 48px; font-weight: 800; color: var(--blue);
    line-height: 1; white-space: nowrap;
  }
  .overall-sub { font-size: 13px; color: var(--text-secondary); text-align: right; }

  /* ── STATUS CARDS ── */
  .status-row {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
    gap: 10px;
    padding: 16px 32px;
    background: var(--gray);
    border-bottom: 1px solid var(--border);
  }
  .status-card {
    background: #fff;
    border-radius: var(--radius);
    padding: 14px 16px;
    border: 1px solid var(--border);
  }
  .status-card-label {
    font-size: 11px; font-weight: 600; text-transform: uppercase;
    letter-spacing: .06em; color: var(--text-secondary); margin-bottom: 6px;
  }
  .status-card-value { font-size: 14px; font-weight: 500; color: var(--dark); line-height: 1.4; }
  .status-card.highlight { border-color: var(--blue); background: var(--blue-light); }
  .status-card.highlight .status-card-value { color: var(--blue); font-weight: 700; }

  /* ── PHASES GRID ── */
  .phases-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(min(100%, 360px), 1fr));
    gap: 16px;
    padding: 20px 32px;
    max-width: 1600px;
    margin: 0 auto;
  }

  /* ── PHASE CARD ── */
  .phase-card {
    background: #fff;
    border-radius: var(--radius);
    border: 1px solid var(--border);
    overflow: hidden;
    border-top: 4px solid var(--phase-color);
  }
  .phase-header {
    padding: 16px 18px 12px;
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 12px;
  }
  .phase-meta { display: flex; gap: 12px; align-items: flex-start; }
  .phase-id {
    font-size: 11px; font-weight: 700; padding: 3px 8px;
    border-radius: 6px; white-space: nowrap;
    background: var(--phase-bg); color: var(--phase-color);
  }
  .phase-title { font-size: 15px; font-weight: 700; color: var(--dark); line-height: 1.3; }
  .phase-duration { font-size: 12px; color: var(--text-secondary); margin-top: 3px; }
  .phase-right { display: flex; flex-direction: column; align-items: flex-end; gap: 4px; flex-shrink: 0; }
  .phase-pct { font-size: 22px; font-weight: 800; color: var(--phase-color); }

  /* badges */
  .badge {
    font-size: 10px; font-weight: 700; padding: 2px 8px;
    border-radius: 99px; text-transform: uppercase; letter-spacing: .04em;
  }
  .badge-done   { background: var(--green-bg); color: var(--green); }
  .badge-active { background: var(--blue-light); color: var(--blue); }
  .badge-todo   { background: #F1F5F9; color: #64748B; }

  /* progress bar */
  .phase-progress-bar {
    height: 5px; background: var(--phase-bg);
    margin: 0 18px 0;
  }
  .phase-progress-fill {
    height: 100%; background: var(--phase-color);
    transition: width .5s ease;
  }

  /* body */
  .phase-body {
    padding: 14px 18px;
    max-height: 380px;
    overflow-y: auto;
  }
  .phase-body::-webkit-scrollbar { width: 4px; }
  .phase-body::-webkit-scrollbar-thumb { background: var(--border); border-radius: 2px; }

  /* subphase */
  .subphase { margin-bottom: 12px; }
  .subphase-header {
    display: flex; justify-content: space-between; align-items: center;
    margin-bottom: 6px; padding-bottom: 4px;
    border-bottom: 1px solid var(--border);
  }
  .subphase-title { font-size: 12px; font-weight: 700; color: var(--blue); text-transform: uppercase; letter-spacing: .04em; }
  .subphase-count { font-size: 11px; color: var(--text-secondary); }

  /* task */
  .task {
    display: flex; align-items: flex-start; gap: 8px;
    padding: 4px 0; font-size: 13px; color: var(--dark); line-height: 1.4;
  }
  .task-done { color: var(--text-secondary); }
  .task-done span:last-child { text-decoration: line-through; }
  .task-icon { font-size: 12px; flex-shrink: 0; margin-top: 1px; }
  .task-icon.done { color: var(--green); }
  .task-icon.todo { color: #CBD5E1; }

  /* footer */
  .phase-footer {
    padding: 8px 18px 12px;
    font-size: 11px; color: var(--text-secondary);
    border-top: 1px solid var(--border); margin-top: 4px;
  }

  /* ── RESPONSIVE ── */

  /* Планшет */
  @media (max-width: 768px) {
    .site-header { padding: 16px 20px; gap: 10px; }
    .site-title  { font-size: 17px; }
    .overall     { padding: 20px; gap: 16px; }
    .status-row  { padding: 12px 20px; gap: 8px;
                   grid-template-columns: repeat(2, 1fr); }
    .phases-grid { padding: 16px 20px; gap: 12px; }
  }

  /* Мобильный */
  @media (max-width: 480px) {
    .site-header { padding: 12px 16px; }
    .site-title  { font-size: 15px; }
    .site-sub    { font-size: 11px; }
    .updated     { font-size: 11px; }

    .overall {
      padding: 14px 16px;
      grid-template-columns: 1fr auto;
      gap: 12px;
    }
    .overall-label { font-size: 12px; }
    .overall-pct   { font-size: 28px; }
    .overall-sub   { font-size: 11px; }

    .status-row {
      padding: 10px 16px;
      grid-template-columns: 1fr 1fr;
      gap: 8px;
    }
    .status-card       { padding: 10px 12px; }
    .status-card-label { font-size: 10px; }
    .status-card-value { font-size: 13px; }

    .phases-grid {
      padding: 12px 16px;
      grid-template-columns: 1fr;
      gap: 10px;
    }

    .phase-header  { padding: 12px 14px 10px; gap: 8px; }
    .phase-meta    { gap: 8px; }
    .phase-title   { font-size: 13px; }
    .phase-duration{ font-size: 11px; }
    .phase-pct     { font-size: 18px; }
    .phase-id      { font-size: 10px; padding: 2px 6px; }

    .phase-progress-bar { margin: 0 14px; }
    .phase-body    { padding: 10px 14px; max-height: none; }
    .phase-footer  { padding: 6px 14px 10px; }

    .task          { font-size: 12px; padding: 3px 0; gap: 6px; }
    .subphase-title{ font-size: 11px; }

    .badge { font-size: 9px; padding: 2px 6px; }
  }

  /* Очень маленький (≤360px, старые Android) */
  @media (max-width: 360px) {
    .status-row { grid-template-columns: 1fr; }
    .overall    { grid-template-columns: 1fr; }
    .overall-pct{ text-align: left; }
  }
</style>
</head>
<body>

<header class="site-header">
  <div>
    <div class="site-title">ТОО «Первое кредитное товарищество»</div>
    <div class="site-sub">Roadmap разработки веб-платформы и мобильного приложения</div>
  </div>
  <div class="updated">Обновлено: ${escHtml(now)}</div>
</header>

<div class="overall">
  <div class="overall-bar-wrap">
    <div class="overall-label">Общий прогресс — ${doneTasks} из ${totalTasks} задач выполнено</div>
    <div class="overall-bar">
      <div class="overall-fill" style="width:${overallPct}%"></div>
    </div>
  </div>
  <div>
    <div class="overall-pct">${overallPct}%</div>
    <div class="overall-sub">${doneTasks} / ${totalTasks}</div>
  </div>
</div>

<div class="status-row">
  <div class="status-card highlight">
    <div class="status-card-label">Текущая фаза</div>
    <div class="status-card-value">${escHtml(currentPhaseName)}</div>
  </div>
  <div class="status-card">
    <div class="status-card-label">Статус</div>
    <div class="status-card-value">${escHtml(status)}</div>
  </div>
  <div class="status-card">
    <div class="status-card-label">Последнее сделанное</div>
    <div class="status-card-value">${escHtml(lastDone)}</div>
  </div>
  <div class="status-card">
    <div class="status-card-label">Следующий шаг</div>
    <div class="status-card-value">${escHtml(nextStep)}</div>
  </div>
  ${blocked !== '—' ? `
  <div class="status-card" style="border-color:#FCA5A5;background:#FEF2F2">
    <div class="status-card-label" style="color:#991B1B">⚠ Заблокировано</div>
    <div class="status-card-value" style="color:#991B1B">${escHtml(blocked)}</div>
  </div>` : ''}
</div>

<div class="phases-grid">
  ${phases.map((p, i) => renderPhase(p, i)).join('\n')}
</div>

</body>
</html>`;

fs.writeFileSync(outFile, html, 'utf8');
console.log(`✓ Dashboard generated → ${outFile}`);
console.log(`  Overall: ${overallPct}% (${doneTasks}/${totalTasks} tasks)`);
console.log(`  Phases parsed: ${phases.length}`);
