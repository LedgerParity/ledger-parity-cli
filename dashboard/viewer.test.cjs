// No browser packages required. Run with: node --test dashboard/viewer.test.cjs
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const path = require('node:path');
const { spawnSync } = require('node:child_process');
const html = fs.readFileSync(path.join(__dirname, 'index.html'), 'utf8');
const source = html.match(/<script>([\s\S]*?)<\/script>/)[1];

function viewer() {
  const elements = new Map();
  const readers = [];
  const context = vm.createContext({
    document: { querySelector(selector) {
      if (!elements.has(selector)) {
        const classes = new Set();
        elements.set(selector, {innerHTML:'', textContent:'', value:'', events:{},
          classList:{add:c => classes.add(c), remove:c => classes.delete(c), contains:c => classes.has(c)},
          addEventListener(name, fn) { this.events[name] = fn; }});
      }
      return elements.get(selector);
    } },
    FileReader: class { constructor() { readers.push(this); } readAsText() {} }
  });
  vm.runInContext(source, context);
  return { context, elements, readers, run:code => vm.runInContext(code, context) };
}

let demo;
function report() {
  if (!demo) {
    const executable = process.env.LEDGER_PARITY_BIN;
    const result = spawnSync(executable || 'go', executable ? ['--demo', '--format', 'json', '--out', '-'] : ['run', './cmd/ledger-parity', '--demo', '--format', 'json', '--out', '-'],
      {cwd:path.join(__dirname, '..'), encoding:'utf8'});
    assert.ifError(result.error);
    assert.equal(result.status, executable ? 3 : 1, result.stderr);
    demo = JSON.parse(result.stdout);
  }
  return structuredClone(demo);
}

test('real CLI demo validates and displays coverage, counts and candidate IDs', () => {
  const v = viewer();
  v.context.report = report();
  v.run('render(validateReport(report))');
  assert.match(v.elements.get('#stats').innerHTML, /Internal records/);
  assert.match(v.elements.get('#coverage').innerHTML, /assertion|unproven/);
  assert.match(v.elements.get('#meta').textContent, /Reconciliation window/);
  const candidates = v.context.report.results.flatMap(r => r.candidate_operation_ids || []);
  assert.ok(candidates.length > 0);
  for (const id of candidates) assert.ok(v.elements.get('#tbody').innerHTML.includes(id));
});

test('report-controlled HTML is escaped and on-chain identities remain searchable', () => {
  const v = viewer();
  const r = report();
  const attack = '<img src=x onerror="globalThis.compromised=true">';
  r.results[0].notes = attack;
  r.coverage.reason = attack;
  r.coverage.network = attack;
  v.context.report = r;
  v.run('render(validateReport(report))');
  assert.ok(!v.elements.get('#tbody').innerHTML.includes('<img'));
  assert.ok(v.elements.get('#tbody').innerHTML.includes('&lt;img'));
  assert.ok(!v.elements.get('#coverage').innerHTML.includes('<img'));
  assert.ok(v.elements.get('#meta').textContent.includes(attack));
  v.context.rows = [{status:'UNKNOWN', discrepancy:'UNRESOLVED', on_chain_payment:{account:'orphan-sender', destination:'orphan-recipient', amount:'922337203685.4775807', asset_code:'XLM', operation_id:'1234567890123456789'}}];
  v.run('renderTable(rows, "orphan-sender")');
  assert.match(v.elements.get('#tbody').innerHTML, /orphan-recipient/);
  assert.match(v.elements.get('#tbody').innerHTML, /922337203685\.4775807/);
  v.run('renderTable(rows, "1234567890123456789")');
  assert.match(v.elements.get('#result-count').textContent, /Showing 1 of 1/);
});

test('invalid shape, inconsistent totals and numeric amounts fail explicitly', () => {
  const v = viewer();
  for (const bad of [null, {}, [], {hash:'abc'}, {...report(), total_matched:999}, {...report(), results:[{status:'SURPRISE'}]}]) {
    v.context.report = bad;
    assert.throws(() => v.run('validateReport(report)'));
  }
  const r = report();
  r.results[0].status = ['EXACT'];
  v.context.report = r;
  assert.throws(() => v.run('validateReport(report)'), /Invalid result status/);
  r.results[0].status = 'EXACT';
  r.results[0].internal_payment.amount = 9.1;
  v.context.report = r;
  assert.throws(() => v.run('validateReport(report)'), /must be text/);
});

test('loading errors hide previous results and stale reads cannot replace a new report', () => {
  const v = viewer();
  v.run('handleFile({name:"first.json", size:1}); handleFile({name:"second.json", size:1})');
  const r = report();
  v.readers[1].onload({target:{result:JSON.stringify(r)}});
  assert.match(v.elements.get('#load-status').textContent, /second.json/);
  v.readers[0].onload({target:{result:'invalid'}});
  assert.match(v.elements.get('#load-status').textContent, /second.json/);
  v.run('handleFile({name:"too-big.json", size:MAX_BYTES+1})');
  assert.ok(v.elements.get('#dashboard').classList.contains('hidden'));
  assert.match(v.elements.get('#load-status').textContent, /32 MiB/);
  v.run('handleFile({name:"bad.json", size:1})');
  v.readers[2].onload({target:{result:'{}'}});
  assert.match(v.elements.get('#load-status').textContent, /Cannot load report/);
  v.run('handleFile({name:"unreadable.json", size:1})');
  v.readers[3].onerror();
  assert.match(v.elements.get('#load-status').textContent, /Cannot read/);
});

test('large result lists are capped without excluding records from search', () => {
  const v = viewer();
  v.context.rows = Array.from({length:501}, (_, i) => ({status:'UNKNOWN', discrepancy:'UNRESOLVED', notes:'record-' + i}));
  v.run('renderTable(rows)');
  assert.equal((v.elements.get('#tbody').innerHTML.match(/<tr>/g) || []).length, 500);
  assert.match(v.elements.get('#result-count').textContent, /Showing 500 of 501/);
  v.run('renderTable(rows, "record-500")');
  assert.match(v.elements.get('#tbody').innerHTML, /record-500/);
  v.run('renderTable(rows, "nothing-matches")');
  assert.match(v.elements.get('#tbody').innerHTML, /No matching results/);
});
