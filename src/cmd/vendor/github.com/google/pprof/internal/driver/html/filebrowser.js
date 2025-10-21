function fbInit(opts) {
  let base = opts.base || '';
  let source = opts.source || '';
  let selectedDir = '';
  let expandedDirs = {};
  const treeEl = document.getElementById('fb-tree');
  const pathEl = document.getElementById('fb-path');
  const baseEl = document.getElementById('fb-base');
  const sourceEl = document.getElementById('fb-source');
  const btnReset = document.getElementById('fb-reset');
  const btnWork = document.getElementById('fb-work');
  const btnSelectDir = document.getElementById('fb-select-dir');

  // Load saved state on initialization
  function loadState() {
    fetch(makeURL('api/getstate'))
      .then(r => r.json())
      .then(state => {
        selectedDir = state.selectedDir || '';
        expandedDirs = state.expandedDirs || {};
        base = state.baseProfile || '';
        source = state.sourceProfile || '';
        
        if (selectedDir) {
          fetchTree(selectedDir);
        } else {
          fetchTree('/');
        }
        updateHeader();
      })
      .catch(() => {
        fetchTree('/');
        updateHeader();
      });
  }

  // Save directory selection
  function saveSelectedDir(path) {
    selectedDir = path;
    fetch(makeURL('api/setdir', { path: path }))
      .catch(err => console.error('Failed to save directory:', err));
  }

  // Save expanded state
  function saveExpandedState(path, expanded) {
    expandedDirs[path] = expanded;
    fetch(makeURL('api/setexpanded', { path: path, expanded: expanded }))
      .catch(err => console.error('Failed to save expanded state:', err));
  }

  function updateHeader() {
    baseEl.textContent = base || 'Not selected';
    sourceEl.textContent = source || 'Not selected';
    // Work button is enabled when base is selected (source is optional)
    btnWork.disabled = !base;
    const tb = document.getElementById('fb-tab-base');
    const ts = document.getElementById('fb-tab-source');
    if (tb) {
      tb.classList.remove('empty', 'filled', 'selected-base');
      if (base) {
        tb.classList.add('selected-base');
      } else {
        tb.classList.add('empty');
      }
    }
    if (ts) {
      ts.classList.remove('empty', 'filled', 'selected-source');
      if (source) {
        ts.classList.add('selected-source');
      } else {
        ts.classList.add('empty');
      }
    }
  }

  function updateTreeMarkers() {
    // Remove all selection markers
    document.querySelectorAll('.fb-row').forEach(row => {
      row.classList.remove('selected-base', 'selected-source');
    });
    // Compare using dataset to avoid CSS selector escaping issues on Windows paths
    if (base) {
      const rows = document.querySelectorAll('.fb-row');
      rows.forEach(r => { if (r.dataset.path === base) r.classList.add('selected-base'); });
    }
    if (source) {
      const rows = document.querySelectorAll('.fb-row');
      rows.forEach(r => { if (r.dataset.path === source) r.classList.add('selected-source'); });
    }
  }

  function fetchTree(path, forceRefresh = false) {
    pathEl.textContent = path || '/';
    treeEl.innerHTML = '<div class="fb-row">Loading…</div>';
    
    const params = { path: path || '/' };
    if (forceRefresh) {
      params.force = 'true';
    }
    
    fetch(makeURL('api/files', params))
      .then(r => r.json())
      .then(data => {
        if (data.error) {
          treeEl.innerHTML = '<div class="fb-row">Error: ' + data.error + '</div>';
          return;
        }
        treeEl.innerHTML = '';
        const frag = document.createDocumentFragment();
        renderLevel(frag, data.items, path || '/');
        treeEl.appendChild(frag);
        updateTreeMarkers();
        
        // Save the selected directory
        if (path && path !== '/') {
          saveSelectedDir(path);
        }
      })
      .catch(err => {
        treeEl.innerHTML = '<div class="fb-row">Error: ' + err + '</div>';
      });
  }

  function renderLevel(parent, items, basePath) {
    if (!items) return;
    for (const item of items) {
      const row = document.createElement('div');
      row.className = 'fb-row';
      row.dataset.path = item.path;

      const twist = document.createElement('span');
      twist.className = 'twisty' + (item.isDir ? '' : ' hidden');
      twist.textContent = item.isDir ? '▸' : '';
      row.appendChild(twist);

      const kind = document.createElement('span');
      kind.className = 'kind ' + (item.isDir ? 'dir' : 'file');
      row.appendChild(kind);

      const name = document.createElement('span');
      name.className = 'name';
      name.textContent = item.name;
      row.appendChild(name);

      if (!item.isDir) {
        const fileActions = document.createElement('span');
        fileActions.className = 'fb-file-actions';
        const b1 = document.createElement('span');
        b1.className = 'link';
        b1.textContent = 'Select base';
        b1.addEventListener('click', (e) => {
          e.stopPropagation();
          setBase(item.path);
        });
        const b2 = document.createElement('span');
        b2.className = 'link';
        b2.textContent = 'Select source';
        b2.addEventListener('click', (e) => {
          e.stopPropagation();
          setSource(item.path);
        });
        fileActions.appendChild(b1);
        fileActions.appendChild(document.createTextNode(' · '));
        fileActions.appendChild(b2);
        row.appendChild(fileActions);
      }

      parent.appendChild(row);

      // Prevent double-click text selection causing highlight bars
      row.addEventListener('mousedown', (e) => {
        if (e.detail > 1) e.preventDefault();
      });

      if (item.isDir) {
        const children = document.createElement('div');
        children.className = 'fb-children hidden';
        parent.appendChild(children);
        let loaded = false;
        
        // Check if this directory should be expanded based on saved state
        const shouldExpand = expandedDirs[item.path] || false;
        if (shouldExpand) {
          twist.textContent = '▾';
          children.classList.remove('hidden');
          // Load children immediately if should be expanded
          loadChildren(item.path, children, twist);
          loaded = true;
        }
        
        row.addEventListener('click', (e) => {
          // Only toggle on row background/name/twisty, ignore clicks on action buttons
          if (e.target && (e.target.classList.contains('link') || e.target.closest('.fb-file-actions'))) {
            return;
          }
          if (!loaded) {
            loadChildren(item.path, children, twist);
            loaded = true;
          } else {
            const hidden = children.classList.toggle('hidden');
            twist.textContent = hidden ? '▸' : '▾';
            saveExpandedState(item.path, !hidden);
          }
        });
      }
    }
  }

  function loadChildren(path, childrenEl, twistEl) {
    childrenEl.innerHTML = '<div class="fb-row">Loading…</div>';
    fetch(makeURL('api/files', { path: path }))
      .then(r => r.json())
      .then(data => {
        childrenEl.innerHTML = '';
        renderLevel(childrenEl, data.items, path);
        childrenEl.classList.remove('hidden');
        twistEl.textContent = '▾';
        updateTreeMarkers();
        saveExpandedState(path, true);
      })
      .catch(() => { 
        childrenEl.innerHTML = '<div class="fb-row">Error</div>'; 
        saveExpandedState(path, false);
      });
  }

  function setBase(p) {
    fetch(makeURL('api/setbase', { path: p }))
      .then(r => { 
        if (r.ok) { 
          base = p; 
          updateHeader(); 
          updateTreeMarkers(); 
        }
      });
  }
  function setSource(p) {
    fetch(makeURL('api/setsource', { path: p }))
      .then(r => { 
        if (r.ok) { 
          source = p; 
          updateHeader(); 
          updateTreeMarkers(); 
        }
      });
  }
  function reset() {
    fetch('api/reset').then(r => { if (r.ok) { base = ''; source = ''; updateHeader(); updateTreeMarkers(); }});
  }
  function work() {
    fetch('api/work').then(r => {
      if (r.ok) {
        // Navigate to graph view after successful work
        window.location.href = '/ui/';
      } else {
        alert('Work failed');
      }
    }).catch(() => alert('Work failed'));
  }
  function selectDirectory() {
    const current = pathEl.textContent || '/';
    const p = prompt('Enter directory path to scan', current);
    if (p) {
      fetchTree(p, true); // Force refresh when manually selecting directory
    }
  }

  function makeURL(path, params) {
    const u = new URL(path, document.URL);
    if (params) for (const [k, v] of Object.entries(params)) u.searchParams.set(k, v);
    return u.toString();
  }

  btnReset.addEventListener('click', reset);
  btnWork.addEventListener('click', work);
  btnSelectDir.addEventListener('click', selectDirectory);
  
  // Load state and initialize
  loadState();

  // Enable top menus and config dialogs on this page
  if (typeof initMenus === 'function') initMenus();
  if (typeof initConfigManager === 'function') initConfigManager();
}