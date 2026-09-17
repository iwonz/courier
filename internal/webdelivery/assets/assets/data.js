(function(){let e=document.createElement(`link`).relList;if(e&&e.supports&&e.supports(`modulepreload`))return;for(let e of document.querySelectorAll(`link[rel="modulepreload"]`))n(e);new MutationObserver(e=>{for(let t of e)if(t.type===`childList`)for(let e of t.addedNodes)e.tagName===`LINK`&&e.rel===`modulepreload`&&n(e)}).observe(document,{childList:!0,subtree:!0});function t(e){let t={};return e.integrity&&(t.integrity=e.integrity),e.referrerPolicy&&(t.referrerPolicy=e.referrerPolicy),t.credentials=e.crossOrigin===`use-credentials`?`include`:e.crossOrigin===`anonymous`?`omit`:`same-origin`,t}function n(e){if(e.ep)return;e.ep=!0;let n=t(e);fetch(e.href,n)}})();var e=globalThis,t=e.ShadowRoot&&(e.ShadyCSS===void 0||e.ShadyCSS.nativeShadow)&&`adoptedStyleSheets`in Document.prototype&&`replace`in CSSStyleSheet.prototype,n=Symbol(),r=new WeakMap,i=class{constructor(e,t,r){if(this._$cssResult$=!0,r!==n)throw Error("CSSResult is not constructable. Use `unsafeCSS` or `css` instead.");this.cssText=e,this.t=t}get styleSheet(){let e=this.o,n=this.t;if(t&&e===void 0){let t=n!==void 0&&n.length===1;t&&(e=r.get(n)),e===void 0&&((this.o=e=new CSSStyleSheet).replaceSync(this.cssText),t&&r.set(n,e))}return e}toString(){return this.cssText}},a=e=>new i(typeof e==`string`?e:e+``,void 0,n),o=(e,...t)=>new i(e.length===1?e[0]:t.reduce((t,n,r)=>t+(e=>{if(!0===e._$cssResult$)return e.cssText;if(typeof e==`number`)return e;throw Error(`Value passed to 'css' function must be a 'css' function result: `+e+`. Use 'unsafeCSS' to pass non-literal values, but take care to ensure page security.`)})(n)+e[r+1],e[0]),e,n),s=(n,r)=>{if(t)n.adoptedStyleSheets=r.map(e=>e instanceof CSSStyleSheet?e:e.styleSheet);else for(let t of r){let r=document.createElement(`style`),i=e.litNonce;i!==void 0&&r.setAttribute(`nonce`,i),r.textContent=t.cssText,n.appendChild(r)}},c=t?e=>e:e=>e instanceof CSSStyleSheet?(e=>{let t=``;for(let n of e.cssRules)t+=n.cssText;return a(t)})(e):e,{is:l,defineProperty:u,getOwnPropertyDescriptor:d,getOwnPropertyNames:ee,getOwnPropertySymbols:te,getPrototypeOf:ne}=Object,f=globalThis,re=f.trustedTypes,ie=re?re.emptyScript:``,ae=f.reactiveElementPolyfillSupport,p=(e,t)=>e,m={toAttribute(e,t){switch(t){case Boolean:e=e?ie:null;break;case Object:case Array:e=e==null?e:JSON.stringify(e)}return e},fromAttribute(e,t){let n=e;switch(t){case Boolean:n=e!==null;break;case Number:n=e===null?null:Number(e);break;case Object:case Array:try{n=JSON.parse(e)}catch{n=null}}return n}},oe=(e,t)=>!l(e,t),h={attribute:!0,type:String,converter:m,reflect:!1,useDefault:!1,hasChanged:oe};Symbol.metadata??=Symbol(`metadata`),f.litPropertyMetadata??=new WeakMap;var g=class extends HTMLElement{static addInitializer(e){this._$Ei(),(this.l??=[]).push(e)}static get observedAttributes(){return this.finalize(),this._$Eh&&[...this._$Eh.keys()]}static createProperty(e,t=h){if(t.state&&(t.attribute=!1),this._$Ei(),this.prototype.hasOwnProperty(e)&&((t=Object.create(t)).wrapped=!0),this.elementProperties.set(e,t),!t.noAccessor){let n=Symbol(),r=this.getPropertyDescriptor(e,n,t);r!==void 0&&u(this.prototype,e,r)}}static getPropertyDescriptor(e,t,n){let{get:r,set:i}=d(this.prototype,e)??{get(){return this[t]},set(e){this[t]=e}};return{get:r,set(t){let a=r?.call(this);i?.call(this,t),this.requestUpdate(e,a,n)},configurable:!0,enumerable:!0}}static getPropertyOptions(e){return this.elementProperties.get(e)??h}static _$Ei(){if(this.hasOwnProperty(p(`elementProperties`)))return;let e=ne(this);e.finalize(),e.l!==void 0&&(this.l=[...e.l]),this.elementProperties=new Map(e.elementProperties)}static finalize(){if(this.hasOwnProperty(p(`finalized`)))return;if(this.finalized=!0,this._$Ei(),this.hasOwnProperty(p(`properties`))){let e=this.properties,t=[...ee(e),...te(e)];for(let n of t)this.createProperty(n,e[n])}let e=this[Symbol.metadata];if(e!==null){let t=litPropertyMetadata.get(e);if(t!==void 0)for(let[e,n]of t)this.elementProperties.set(e,n)}this._$Eh=new Map;for(let[e,t]of this.elementProperties){let n=this._$Eu(e,t);n!==void 0&&this._$Eh.set(n,e)}this.elementStyles=this.finalizeStyles(this.styles)}static finalizeStyles(e){let t=[];if(Array.isArray(e)){let n=new Set(e.flat(1/0).reverse());for(let e of n)t.unshift(c(e))}else e!==void 0&&t.push(c(e));return t}static _$Eu(e,t){let n=t.attribute;return!1===n?void 0:typeof n==`string`?n:typeof e==`string`?e.toLowerCase():void 0}constructor(){super(),this._$Ep=void 0,this.isUpdatePending=!1,this.hasUpdated=!1,this._$Em=null,this._$Ev()}_$Ev(){this._$ES=new Promise(e=>this.enableUpdating=e),this._$AL=new Map,this._$E_(),this.requestUpdate(),this.constructor.l?.forEach(e=>e(this))}addController(e){(this._$EO??=new Set).add(e),this.renderRoot!==void 0&&this.isConnected&&e.hostConnected?.()}removeController(e){this._$EO?.delete(e)}_$E_(){let e=new Map,t=this.constructor.elementProperties;for(let n of t.keys())this.hasOwnProperty(n)&&(e.set(n,this[n]),delete this[n]);e.size>0&&(this._$Ep=e)}createRenderRoot(){let e=this.shadowRoot??this.attachShadow(this.constructor.shadowRootOptions);return s(e,this.constructor.elementStyles),e}connectedCallback(){this.renderRoot??=this.createRenderRoot(),this.enableUpdating(!0),this._$EO?.forEach(e=>e.hostConnected?.())}enableUpdating(e){}disconnectedCallback(){this._$EO?.forEach(e=>e.hostDisconnected?.())}attributeChangedCallback(e,t,n){this._$AK(e,n)}_$ET(e,t){let n=this.constructor.elementProperties.get(e),r=this.constructor._$Eu(e,n);if(r!==void 0&&!0===n.reflect){let i=(n.converter?.toAttribute===void 0?m:n.converter).toAttribute(t,n.type);this._$Em=e,i==null?this.removeAttribute(r):this.setAttribute(r,i),this._$Em=null}}_$AK(e,t){let n=this.constructor,r=n._$Eh.get(e);if(r!==void 0&&this._$Em!==r){let e=n.getPropertyOptions(r),i=typeof e.converter==`function`?{fromAttribute:e.converter}:e.converter?.fromAttribute===void 0?m:e.converter;this._$Em=r;let a=i.fromAttribute(t,e.type);this[r]=a??this._$Ej?.get(r)??a,this._$Em=null}}requestUpdate(e,t,n,r=!1,i){if(e!==void 0){let a=this.constructor;if(!1===r&&(i=this[e]),n??=a.getPropertyOptions(e),!((n.hasChanged??oe)(i,t)||n.useDefault&&n.reflect&&i===this._$Ej?.get(e)&&!this.hasAttribute(a._$Eu(e,n))))return;this.C(e,t,n)}!1===this.isUpdatePending&&(this._$ES=this._$EP())}C(e,t,{useDefault:n,reflect:r,wrapped:i},a){n&&!(this._$Ej??=new Map).has(e)&&(this._$Ej.set(e,a??t??this[e]),!0!==i||a!==void 0)||(this._$AL.has(e)||(this.hasUpdated||n||(t=void 0),this._$AL.set(e,t)),!0===r&&this._$Em!==e&&(this._$Eq??=new Set).add(e))}async _$EP(){this.isUpdatePending=!0;try{await this._$ES}catch(e){Promise.reject(e)}let e=this.scheduleUpdate();return e!=null&&await e,!this.isUpdatePending}scheduleUpdate(){return this.performUpdate()}performUpdate(){if(!this.isUpdatePending)return;if(!this.hasUpdated){if(this.renderRoot??=this.createRenderRoot(),this._$Ep){for(let[e,t]of this._$Ep)this[e]=t;this._$Ep=void 0}let e=this.constructor.elementProperties;if(e.size>0)for(let[t,n]of e){let{wrapped:e}=n,r=this[t];!0!==e||this._$AL.has(t)||r===void 0||this.C(t,void 0,n,r)}}let e=!1,t=this._$AL;try{e=this.shouldUpdate(t),e?(this.willUpdate(t),this._$EO?.forEach(e=>e.hostUpdate?.()),this.update(t)):this._$EM()}catch(t){throw e=!1,this._$EM(),t}e&&this._$AE(t)}willUpdate(e){}_$AE(e){this._$EO?.forEach(e=>e.hostUpdated?.()),this.hasUpdated||(this.hasUpdated=!0,this.firstUpdated(e)),this.updated(e)}_$EM(){this._$AL=new Map,this.isUpdatePending=!1}get updateComplete(){return this.getUpdateComplete()}getUpdateComplete(){return this._$ES}shouldUpdate(e){return!0}update(e){this._$Eq&&=this._$Eq.forEach(e=>this._$ET(e,this[e])),this._$EM()}updated(e){}firstUpdated(e){}};g.elementStyles=[],g.shadowRootOptions={mode:`open`},g[p(`elementProperties`)]=new Map,g[p(`finalized`)]=new Map,ae?.({ReactiveElement:g}),(f.reactiveElementVersions??=[]).push(`2.1.2`);var _=globalThis,v=e=>e,y=_.trustedTypes,b=y?y.createPolicy(`lit-html`,{createHTML:e=>e}):void 0,x=`$lit$`,S=`lit$${Math.random().toFixed(9).slice(2)}$`,se=`?`+S,ce=`<${se}>`,C=document,w=()=>C.createComment(``),T=e=>e===null||typeof e!=`object`&&typeof e!=`function`,E=Array.isArray,le=e=>E(e)||typeof e?.[Symbol.iterator]==`function`,D=`[ 	
\f\r]`,O=/<(?:(!--|\/[^a-zA-Z])|(\/?[a-zA-Z][^>\s]*)|(\/?$))/g,ue=/-->/g,k=/>/g,A=RegExp(`>|${D}(?:([^\\s"'>=/]+)(${D}*=${D}*(?:[^ \t\n\f\r"'\`<>=]|("|')|))|$)`,`g`),de=/'/g,fe=/"/g,j=/^(?:script|style|textarea|title)$/i,M=(e=>(t,...n)=>({_$litType$:e,strings:t,values:n}))(1),N=Symbol.for(`lit-noChange`),P=Symbol.for(`lit-nothing`),F=new WeakMap,I=C.createTreeWalker(C,129);function L(e,t){if(!E(e)||!e.hasOwnProperty(`raw`))throw Error(`invalid template strings array`);return b===void 0?t:b.createHTML(t)}var pe=(e,t)=>{let n=e.length-1,r=[],i,a=t===2?`<svg>`:t===3?`<math>`:``,o=O;for(let t=0;t<n;t++){let n=e[t],s,c,l=-1,u=0;for(;u<n.length&&(o.lastIndex=u,c=o.exec(n),c!==null);)u=o.lastIndex,o===O?c[1]===`!--`?o=ue:c[1]===void 0?c[2]===void 0?c[3]!==void 0&&(o=A):(j.test(c[2])&&(i=RegExp(`</`+c[2],`g`)),o=A):o=k:o===A?c[0]===`>`?(o=i??O,l=-1):c[1]===void 0?l=-2:(l=o.lastIndex-c[2].length,s=c[1],o=c[3]===void 0?A:c[3]===`"`?fe:de):o===fe||o===de?o=A:o===ue||o===k?o=O:(o=A,i=void 0);let d=o===A&&e[t+1].startsWith(`/>`)?` `:``;a+=o===O?n+ce:l>=0?(r.push(s),n.slice(0,l)+x+n.slice(l)+S+d):n+S+(l===-2?t:d)}return[L(e,a+(e[n]||`<?>`)+(t===2?`</svg>`:t===3?`</math>`:``)),r]},R=class e{constructor({strings:t,_$litType$:n},r){let i;this.parts=[];let a=0,o=0,s=t.length-1,c=this.parts,[l,u]=pe(t,n);if(this.el=e.createElement(l,r),I.currentNode=this.el.content,n===2||n===3){let e=this.el.content.firstChild;e.replaceWith(...e.childNodes)}for(;(i=I.nextNode())!==null&&c.length<s;){if(i.nodeType===1){if(i.hasAttributes())for(let e of i.getAttributeNames())if(e.endsWith(x)){let t=u[o++],n=i.getAttribute(e).split(S),r=/([.?@])?(.*)/.exec(t);c.push({type:1,index:a,name:r[2],strings:n,ctor:r[1]===`.`?he:r[1]===`?`?ge:r[1]===`@`?_e:V}),i.removeAttribute(e)}else e.startsWith(S)&&(c.push({type:6,index:a}),i.removeAttribute(e));if(j.test(i.tagName)){let e=i.textContent.split(S),t=e.length-1;if(t>0){i.textContent=y?y.emptyScript:``;for(let n=0;n<t;n++)i.append(e[n],w()),I.nextNode(),c.push({type:2,index:++a});i.append(e[t],w())}}}else if(i.nodeType===8){if(i.data===se)c.push({type:2,index:a});else{let e=-1;for(;(e=i.data.indexOf(S,e+1))!==-1;)c.push({type:7,index:a}),e+=S.length-1}}a++}}static createElement(e,t){let n=C.createElement(`template`);return n.innerHTML=e,n}};function z(e,t,n=e,r){if(t===N)return t;let i=r===void 0?n._$Cl:n._$Co?.[r],a=T(t)?void 0:t._$litDirective$;return i?.constructor!==a&&(i?._$AO?.(!1),a===void 0?i=void 0:(i=new a(e),i._$AT(e,n,r)),r===void 0?n._$Cl=i:(n._$Co??=[])[r]=i),i!==void 0&&(t=z(e,i._$AS(e,t.values),i,r)),t}var me=class{constructor(e,t){this._$AV=[],this._$AN=void 0,this._$AD=e,this._$AM=t}get parentNode(){return this._$AM.parentNode}get _$AU(){return this._$AM._$AU}u(e){let{el:{content:t},parts:n}=this._$AD,r=(e?.creationScope??C).importNode(t,!0);I.currentNode=r;let i=I.nextNode(),a=0,o=0,s=n[0];for(;s!==void 0;){if(a===s.index){let t;s.type===2?t=new B(i,i.nextSibling,this,e):s.type===1?t=new s.ctor(i,s.name,s.strings,this,e):s.type===6&&(t=new ve(i,this,e)),this._$AV.push(t),s=n[++o]}a!==s?.index&&(i=I.nextNode(),a++)}return I.currentNode=C,r}p(e){let t=0;for(let n of this._$AV)n!==void 0&&(n.strings===void 0?n._$AI(e[t]):(n._$AI(e,n,t),t+=n.strings.length-2)),t++}},B=class e{get _$AU(){return this._$AM?._$AU??this._$Cv}constructor(e,t,n,r){this.type=2,this._$AH=P,this._$AN=void 0,this._$AA=e,this._$AB=t,this._$AM=n,this.options=r,this._$Cv=r?.isConnected??!0}get parentNode(){let e=this._$AA.parentNode,t=this._$AM;return t!==void 0&&e?.nodeType===11&&(e=t.parentNode),e}get startNode(){return this._$AA}get endNode(){return this._$AB}_$AI(e,t=this){e=z(this,e,t),T(e)?e===P||e==null||e===``?(this._$AH!==P&&this._$AR(),this._$AH=P):e!==this._$AH&&e!==N&&this._(e):e._$litType$===void 0?e.nodeType===void 0?le(e)?this.k(e):this._(e):this.T(e):this.$(e)}O(e){return this._$AA.parentNode.insertBefore(e,this._$AB)}T(e){this._$AH!==e&&(this._$AR(),this._$AH=this.O(e))}_(e){this._$AH!==P&&T(this._$AH)?this._$AA.nextSibling.data=e:this.T(C.createTextNode(e)),this._$AH=e}$(e){let{values:t,_$litType$:n}=e,r=typeof n==`number`?this._$AC(e):(n.el===void 0&&(n.el=R.createElement(L(n.h,n.h[0]),this.options)),n);if(this._$AH?._$AD===r)this._$AH.p(t);else{let e=new me(r,this),n=e.u(this.options);e.p(t),this.T(n),this._$AH=e}}_$AC(e){let t=F.get(e.strings);return t===void 0&&F.set(e.strings,t=new R(e)),t}k(t){E(this._$AH)||(this._$AH=[],this._$AR());let n=this._$AH,r,i=0;for(let a of t)i===n.length?n.push(r=new e(this.O(w()),this.O(w()),this,this.options)):r=n[i],r._$AI(a),i++;i<n.length&&(this._$AR(r&&r._$AB.nextSibling,i),n.length=i)}_$AR(e=this._$AA.nextSibling,t){for(this._$AP?.(!1,!0,t);e!==this._$AB;){let t=v(e).nextSibling;v(e).remove(),e=t}}setConnected(e){this._$AM===void 0&&(this._$Cv=e,this._$AP?.(e))}},V=class{get tagName(){return this.element.tagName}get _$AU(){return this._$AM._$AU}constructor(e,t,n,r,i){this.type=1,this._$AH=P,this._$AN=void 0,this.element=e,this.name=t,this._$AM=r,this.options=i,n.length>2||n[0]!==``||n[1]!==``?(this._$AH=Array(n.length-1).fill(new String),this.strings=n):this._$AH=P}_$AI(e,t=this,n,r){let i=this.strings,a=!1;if(i===void 0)e=z(this,e,t,0),a=!T(e)||e!==this._$AH&&e!==N,a&&(this._$AH=e);else{let r=e,o,s;for(e=i[0],o=0;o<i.length-1;o++)s=z(this,r[n+o],t,o),s===N&&(s=this._$AH[o]),a||=!T(s)||s!==this._$AH[o],s===P?e=P:e!==P&&(e+=(s??``)+i[o+1]),this._$AH[o]=s}a&&!r&&this.j(e)}j(e){e===P?this.element.removeAttribute(this.name):this.element.setAttribute(this.name,e??``)}},he=class extends V{constructor(){super(...arguments),this.type=3}j(e){this.element[this.name]=e===P?void 0:e}},ge=class extends V{constructor(){super(...arguments),this.type=4}j(e){this.element.toggleAttribute(this.name,!!e&&e!==P)}},_e=class extends V{constructor(e,t,n,r,i){super(e,t,n,r,i),this.type=5}_$AI(e,t=this){if((e=z(this,e,t,0)??P)===N)return;let n=this._$AH,r=e===P&&n!==P||e.capture!==n.capture||e.once!==n.once||e.passive!==n.passive,i=e!==P&&(n===P||r);r&&this.element.removeEventListener(this.name,this,n),i&&this.element.addEventListener(this.name,this,e),this._$AH=e}handleEvent(e){typeof this._$AH==`function`?this._$AH.call(this.options?.host??this.element,e):this._$AH.handleEvent(e)}},ve=class{constructor(e,t,n){this.element=e,this.type=6,this._$AN=void 0,this._$AM=t,this.options=n}get _$AU(){return this._$AM._$AU}_$AI(e){z(this,e)}},ye=_.litHtmlPolyfillSupport;ye?.(R,B),(_.litHtmlVersions??=[]).push(`3.3.3`);var be=(e,t,n)=>{let r=n?.renderBefore??t,i=r._$litPart$;if(i===void 0){let e=n?.renderBefore??null;r._$litPart$=i=new B(t.insertBefore(w(),e),e,void 0,n??{})}return i._$AI(e),i},H=globalThis,U=class extends g{constructor(){super(...arguments),this.renderOptions={host:this},this._$Do=void 0}createRenderRoot(){let e=super.createRenderRoot();return this.renderOptions.renderBefore??=e.firstChild,e}update(e){let t=this.render();this.hasUpdated||(this.renderOptions.isConnected=this.isConnected),super.update(e),this._$Do=be(t,this.renderRoot,this.renderOptions)}connectedCallback(){super.connectedCallback(),this._$Do?.setConnected(!0)}disconnectedCallback(){super.disconnectedCallback(),this._$Do?.setConnected(!1)}render(){return N}};U._$litElement$=!0,U.finalized=!0,H.litElementHydrateSupport?.({LitElement:U});var xe=H.litElementPolyfillSupport;xe?.({LitElement:U}),(H.litElementVersions??=[]).push(`4.2.2`);var Se=o`
  :host {
    color: var(--courier-color-text, #151714);
    font-family: var(--courier-font-sans, sans-serif);
  }

  button,
  select {
    min-height: 2.75rem;
    border: 1px solid var(--courier-color-border, #c8cdbf);
    border-radius: var(--courier-radius-sm, 0.25rem);
    color: inherit;
    background: var(--courier-color-surface-raised, #fff);
    font: inherit;
  }

  button {
    padding: 0.625rem 1rem;
    box-shadow: 0 1px 0 rgb(16 18 15 / 0.08);
    cursor: pointer;
    font-weight: 750;
    letter-spacing: -0.01em;
    transition: background var(--courier-duration, 160ms) var(--courier-ease, ease), border-color var(--courier-duration, 160ms) var(--courier-ease, ease), transform var(--courier-duration, 160ms) var(--courier-ease, ease);
  }

  button:hover:not(:disabled) {
    border-color: var(--courier-color-border-strong, #8e9587);
    background: color-mix(in srgb, var(--courier-color-accent, #d4ff45) 18%, var(--courier-color-surface-raised, #fff));
    transform: translateY(-1px);
  }

  button:active:not(:disabled) {
    transform: translateY(0);
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  select {
    padding: 0.5rem 2rem 0.5rem 0.75rem;
  }

  button:focus-visible,
  select:focus-visible {
    outline: 3px solid var(--courier-color-accent, #ad431d);
    outline-offset: 2px;
  }
`;o`
  label {
    display: grid;
    gap: var(--courier-space-1, 0.25rem);
    color: var(--courier-color-muted, #596054);
    font-family: var(--courier-font-mono, monospace);
    font-size: 0.6875rem;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }
`;var Ce=o`
  input:not([type="checkbox"]):not([type="file"]),
  select {
    appearance: none;
    width: 100%;
    min-width: 0;
    min-height: 2.75rem;
    padding: 0.65rem 0.8rem;
    border: 1px solid var(--courier-color-border-strong, #8e9587);
    border-radius: var(--courier-radius-sm, 0.25rem);
    color: var(--courier-color-text, #151714);
    background-color: var(--courier-color-surface-raised, #fff);
    box-shadow: inset 0 1px 2px rgb(16 18 15 / 0.06);
    font: inherit;
    transition: border-color var(--courier-duration, 160ms) var(--courier-ease, ease), box-shadow var(--courier-duration, 160ms) var(--courier-ease, ease), background var(--courier-duration, 160ms) var(--courier-ease, ease);
  }

  select {
    padding-right: 2.75rem;
    background-image: linear-gradient(45deg, transparent 50%, var(--courier-color-muted, #596054) 50%), linear-gradient(135deg, var(--courier-color-muted, #596054) 50%, transparent 50%);
    background-position: calc(100% - 1.05rem) 50%, calc(100% - 0.72rem) 50%;
    background-repeat: no-repeat;
    background-size: 0.35rem 0.35rem, 0.35rem 0.35rem;
  }

  input:not([type="checkbox"]):not([type="file"]):hover,
  select:hover {
    border-color: var(--courier-color-accent, #d4ff45);
  }

  input:not([type="checkbox"]):not([type="file"]):focus-visible,
  select:focus-visible,
  input[type="checkbox"]:focus-visible {
    outline: 3px solid var(--courier-color-accent, #ad431d);
    outline-offset: 2px;
  }

  input[type="number"] { -moz-appearance: textfield; }
  input[type="number"]::-webkit-inner-spin-button,
  input[type="number"]::-webkit-outer-spin-button { margin: 0; appearance: none; }

  input[type="checkbox"] {
    appearance: none;
    position: relative;
    width: 2.75rem;
    height: 1.55rem;
    margin: 0;
    border: 1px solid var(--courier-color-border-strong, #8e9587);
    border-radius: 999px;
    background: var(--courier-color-field, #e7e9dc);
    box-shadow: inset 0 1px 3px rgb(16 18 15 / 0.12);
    cursor: pointer;
    transition: border-color var(--courier-duration, 160ms) var(--courier-ease, ease), background var(--courier-duration, 160ms) var(--courier-ease, ease);
  }

  input[type="checkbox"]::after {
    content: "";
    position: absolute;
    top: 0.2rem;
    left: 0.2rem;
    width: 1.05rem;
    height: 1.05rem;
    border-radius: 50%;
    background: var(--courier-color-muted, #596054);
    box-shadow: 0 1px 2px rgb(16 18 15 / 0.22);
    transition: transform var(--courier-duration, 160ms) var(--courier-ease, ease), background var(--courier-duration, 160ms) var(--courier-ease, ease);
  }

  input[type="checkbox"]:checked {
    border-color: var(--courier-color-accent, #d4ff45);
    background: color-mix(in srgb, var(--courier-color-accent, #d4ff45) 42%, var(--courier-color-field, #e7e9dc));
  }

  input[type="checkbox"]:checked::after {
    background: var(--courier-color-accent-ink, #151714);
    transform: translateX(1.18rem);
  }

  .courier-file-action {
    position: relative;
    display: inline-flex;
    width: max-content;
    max-width: 100%;
    min-height: 2.75rem;
    align-items: center;
    gap: 0.6rem;
    padding: 0.65rem 0.9rem;
    border: 1px solid var(--courier-color-accent, #d4ff45);
    border-radius: var(--courier-radius-sm, 0.25rem);
    color: var(--courier-color-accent-ink, #151714);
    background: var(--courier-color-accent, #d4ff45);
    box-shadow: 0 1px 0 rgb(16 18 15 / 0.12);
    font-family: var(--courier-font-sans, sans-serif);
    font-size: 0.875rem;
    font-weight: 800;
    letter-spacing: -0.01em;
    cursor: pointer;
  }

  .courier-file-action input[type="file"] {
    position: absolute;
    width: 1px;
    height: 1px;
    overflow: hidden;
    clip-path: inset(50%);
    white-space: nowrap;
  }

  .courier-file-action:has(input[type="file"]:focus-visible) {
    outline: 3px solid var(--courier-color-accent, #ad431d);
    outline-offset: 2px;
  }
`,we=class extends U{constructor(...e){super(...e),this.disabled=!1,this.type=`button`,this.variant=`secondary`}static{this.properties={disabled:{type:Boolean,reflect:!0},type:{type:String,reflect:!0},variant:{type:String,reflect:!0}}}static{this.styles=[Se,o`
    :host { display: inline-flex; }
    button { width: 100%; }
    :host([variant="primary"]) button {
      border-color: var(--courier-color-accent, #d4ff45);
      color: var(--courier-color-accent-ink, #151714);
      background: var(--courier-color-accent, #d4ff45);
      font-weight: 800;
    }
    :host([variant="primary"]) button:hover:not(:disabled) {
      background: color-mix(in srgb, var(--courier-color-accent, #d4ff45) 86%, white);
    }
  `]}forwardSubmit(){this.type===`submit`&&this.closest(`form`)?.requestSubmit()}render(){return M`<button type=${this.type} ?disabled=${this.disabled} @click=${this.forwardSubmit}><slot></slot></button>`}},Te=`data:image/svg+xml,%3csvg%20role='img'%20viewBox='0%200%2024%2024'%20xmlns='http://www.w3.org/2000/svg'%3e%3ctitle%3eAlpine%20Linux%3c/title%3e%3cpath%20d='M5.998%201.607L0%2012l5.998%2010.393h12.004L24%2012%2018.002%201.607H5.998zM9.965%207.12L12.66%209.9l1.598%201.595.002-.002%202.41%202.363c-.2.14-.386.252-.563.344a3.756%203.756%200%2001-.496.217%202.702%202.702%200%2001-.425.111c-.131.023-.25.034-.358.034-.13%200-.242-.014-.338-.034a1.317%201.317%200%2001-.24-.072.95.95%200%2001-.2-.113l-1.062-1.092-3.039-3.041-1.1%201.053-3.07%203.072a.974.974%200%2001-.2.111%201.274%201.274%200%2001-.237.073c-.096.02-.209.033-.338.033-.108%200-.227-.009-.358-.031a2.7%202.7%200%2001-.425-.114%203.748%203.748%200%2001-.496-.217%205.228%205.228%200%2001-.563-.343l6.803-6.727zm4.72.785l4.579%204.598%201.382%201.353a5.24%205.24%200%2001-.564.344%203.73%203.73%200%2001-.494.217%202.697%202.697%200%2001-.426.111c-.13.023-.251.034-.36.034-.129%200-.241-.014-.337-.034a1.285%201.285%200%2001-.385-.146c-.033-.02-.05-.036-.053-.04l-1.232-1.218-2.111-2.111-.334.334L12.79%209.8l1.896-1.897zm-5.966%204.12v2.529a2.128%202.128%200%2001-.356-.035%202.765%202.765%200%2001-.422-.116%203.708%203.708%200%2001-.488-.214%205.217%205.217%200%2001-.555-.34l1.82-1.825Z'/%3e%3c/svg%3e`,Ee=`data:image/svg+xml,%3csvg%20role='img'%20viewBox='0%200%2024%2024'%20xmlns='http://www.w3.org/2000/svg'%3e%3ctitle%3eArch%20Linux%3c/title%3e%3cpath%20d='M11.39.605C10.376%203.092%209.764%204.72%208.635%207.132c.693.734%201.543%201.589%202.923%202.554-1.484-.61-2.496-1.224-3.252-1.86C6.86%2010.842%204.596%2015.138%200%2023.395c3.612-2.085%206.412-3.37%209.021-3.862a6.61%206.61%200%2001-.171-1.547l.003-.115c.058-2.315%201.261-4.095%202.687-3.973%201.426.12%202.534%202.096%202.478%204.409a6.52%206.52%200%2001-.146%201.243c2.58.505%205.352%201.787%208.914%203.844-.702-1.293-1.33-2.459-1.929-3.57-.943-.73-1.926-1.682-3.933-2.713%201.38.359%202.367.772%203.137%201.234-6.09-11.334-6.582-12.84-8.67-17.74zM22.898%2021.36v-.623h-.234v-.084h.562v.084h-.234v.623h.331v-.707h.142l.167.5.034.107a2.26%202.26%200%2001.038-.114l.17-.493H24v.707h-.091v-.593l-.206.593h-.084l-.205-.602v.602h-.091'/%3e%3c/svg%3e`,De=`data:image/svg+xml,%3csvg%20role='img'%20viewBox='0%200%2024%2024'%20xmlns='http://www.w3.org/2000/svg'%3e%3ctitle%3ecurl%3c/title%3e%3cpath%20d='M.803%2014.8169c0-.5342.433-.9665.9665-.9665.5335%200%20.9665.4323.9665.9665%200%20.5335-.433.9657-.9665.9657-.5335%200-.9666-.4322-.9666-.9657m2.736%200c0-.1963-.0532-.376-.1119-.5525-.2344-.7024-.876-1.2169-1.6575-1.2169-.1249%200-.2344.0465-.3524.0708C.6149%2013.2865%200%2013.9646%200%2014.817c0%20.9764.7923%201.7694%201.7695%201.7694.9772%200%201.7694-.793%201.7694-1.7694m-1.7694-7.149c.5335%200%20.9665.433.9665.9665%200%20.5335-.433.9665-.9665.9665-.5343%200-.9666-.433-.9666-.9665%200-.5335.4323-.9665.9666-.9665m0%202.7359c.9772%200%201.7694-.7923%201.7694-1.7694%200-.1956-.0532-.376-.1119-.5525-.2344-.7024-.8767-1.2169-1.6575-1.2169-.1249%200-.2344.0465-.3524.0716C.6149%207.104%200%207.782%200%208.6344c0%20.9771.7923%201.7694%201.7695%201.7694m13.221-5.694c-.5342%200-.9665-.433-.9665-.9664a.966.966%200%2001.9666-.9665c.5335%200%20.9658.4322.9658.9665%200%20.5334-.4323.9664-.9658.9664m-9.6%2016.5133c-.5335%200-.9666-.433-.9666-.9665%200-.5342.433-.9665.9666-.9665a.966.966%200%2001.9665.9665c0%20.5335-.4323.9665-.9665.9665m9.6-19.2491c-.978%200-1.7695.7922-1.7695%201.7694%200%20.2085.0525.4025.1187.5882L5.039%2018.5581c-.803.1681-1.4179.8462-1.4179%201.6985%200%20.9772.7923%201.7694%201.7695%201.7694.9772%200%201.7694-.7922%201.7694-1.7694%200-.1963-.0525-.3759-.111-.5525l8.3427-14.2728c.7778-.1865%201.3683-.8531%201.3683-1.688%200-.977-.793-1.7693-1.7694-1.7693m7.24%202.7359c-.5343%200-.9666-.433-.9666-.9665a.966.966%200%2001.9665-.9665c.5335%200%20.9666.4322.9666.9665%200%20.5334-.433.9665-.9666.9665M12.6313%2021.223c-.5343%200-.9665-.433-.9665-.9665a.966.966%200%2001.9665-.9665c.5335%200%20.9658.4323.9658.9665%200%20.5335-.4323.9665-.9658.9665M22.2305%201.974c-.9772%200-1.7694.7922-1.7694%201.7694%200%20.2085.0525.4025.1187.5882l-8.3009%2014.2265c-.8021.1681-1.417.8462-1.417%201.6985%200%20.9772.7922%201.7694%201.7694%201.7694.9764%200%201.7687-.7922%201.7687-1.7694%200-.1963-.0525-.3759-.1111-.5525l8.3427-14.2728C23.4094%205.2448%2024%204.5782%2024%203.7433c0-.977-.7923-1.7693-1.7695-1.7693'/%3e%3c/svg%3e`,Oe=`data:image/svg+xml,%3csvg%20role='img'%20viewBox='0%200%2024%2024'%20xmlns='http://www.w3.org/2000/svg'%3e%3ctitle%3eDebian%3c/title%3e%3cpath%20d='M13.88%2012.685c-.4%200%20.08.2.601.28.14-.1.27-.22.39-.33a3.001%203.001%200%2001-.99.05m2.14-.53c.23-.33.4-.69.47-1.06-.06.27-.2.5-.33.73-.75.47-.07-.27%200-.56-.8%201.01-.11.6-.14.89m.781-2.05c.05-.721-.14-.501-.2-.221.07.04.13.5.2.22M12.38.31c.2.04.45.07.42.12.23-.05.28-.1-.43-.12m.43.12l-.15.03.14-.01V.43m6.633%209.944c.02.64-.2.95-.38%201.5l-.35.181c-.28.54.03.35-.17.78-.44.39-1.34%201.22-1.62%201.301-.201%200%20.14-.25.19-.34-.591.4-.481.6-1.371.85l-.03-.06c-2.221%201.04-5.303-1.02-5.253-3.842-.03.17-.07.13-.12.2a3.551%203.552%200%20012.001-3.501%203.361%203.362%200%20013.732.48%203.341%203.342%200%2000-2.721-1.3c-1.18.01-2.281.76-2.651%201.57-.6.38-.67%201.47-.93%201.661-.361%202.601.66%203.722%202.38%205.042.27.19.08.21.12.35a4.702%204.702%200%2001-1.53-1.16c.23.33.47.66.8.91-.55-.18-1.27-1.3-1.48-1.35.93%201.66%203.78%202.921%205.261%202.3a6.203%206.203%200%2001-2.33-.28c-.33-.16-.77-.51-.7-.57a5.802%205.803%200%20005.902-.84c.44-.35.93-.94%201.07-.95-.2.32.04.16-.12.44.44-.72-.2-.3.46-1.24l.24.33c-.09-.6.74-1.321.66-2.262.19-.3.2.3%200%20.97.29-.74.08-.85.15-1.46.08.2.18.42.23.63-.18-.7.2-1.2.28-1.6-.09-.05-.28.3-.32-.53%200-.37.1-.2.14-.28-.08-.05-.26-.32-.38-.861.08-.13.22.33.34.34-.08-.42-.2-.75-.2-1.08-.34-.68-.12.1-.4-.3-.34-1.091.3-.25.34-.74.54.77.84%201.96.981%202.46-.1-.6-.28-1.2-.49-1.76.16.07-.26-1.241.21-.37A7.823%207.824%200%200017.702%201.6c.18.17.42.39.33.42-.75-.45-.62-.48-.73-.67-.61-.25-.65.02-1.06%200C15.082.73%2014.862.8%2013.8.4l.05.23c-.77-.25-.9.1-1.73%200-.05-.04.27-.14.53-.18-.741.1-.701-.14-1.431.03.17-.13.36-.21.55-.32-.6.04-1.44.35-1.18.07C9.6.68%207.847%201.3%206.867%202.22L6.838%202c-.45.54-1.96%201.611-2.08%202.311l-.131.03c-.23.4-.38.85-.57%201.261-.3.52-.45.2-.4.28-.6%201.22-.9%202.251-1.16%203.102.18.27%200%201.65.07%202.76-.3%205.463%203.84%2010.776%208.363%2012.006.67.23%201.65.23%202.49.25-.99-.28-1.12-.15-2.08-.49-.7-.32-.85-.7-1.34-1.13l.2.35c-.971-.34-.57-.42-1.361-.67l.21-.27c-.31-.03-.83-.53-.97-.81l-.34.01c-.41-.501-.63-.871-.61-1.161l-.111.2c-.13-.21-1.52-1.901-.8-1.511-.13-.12-.31-.2-.5-.55l.14-.17c-.35-.44-.64-1.02-.62-1.2.2.24.32.3.45.33-.88-2.172-.93-.12-1.601-2.202l.15-.02c-.1-.16-.18-.34-.26-.51l.06-.6c-.63-.74-.18-3.102-.09-4.402.07-.54.53-1.1.88-1.981l-.21-.04c.4-.71%202.341-2.872%203.241-2.761.43-.55-.09%200-.18-.14.96-.991%201.26-.7%201.901-.88.7-.401-.6.16-.27-.151%201.2-.3.85-.7%202.421-.85.16.1-.39.14-.52.26%201-.49%203.151-.37%204.562.27%201.63.77%203.461%203.011%203.531%205.132l.08.02c-.04.85.13%201.821-.17%202.711l.2-.42M9.54%2013.236l-.05.28c.26.35.47.73.8%201.01-.24-.47-.42-.66-.75-1.3m.62-.02c-.14-.15-.22-.34-.31-.52.08.32.26.6.43.88l-.12-.36m10.945-2.382l-.07.15c-.1.76-.34%201.511-.69%202.212.4-.73.65-1.541.75-2.362M12.45.12c.27-.1.66-.05.95-.12-.37.03-.74.05-1.1.1l.15.02M3.006%205.142c.07.57-.43.8.11.42.3-.66-.11-.18-.1-.42m-.64%202.661c.12-.39.15-.62.2-.84-.35.44-.17.53-.2.83'/%3e%3c/svg%3e`,ke=`data:image/svg+xml,%3csvg%20role='img'%20viewBox='0%200%2024%2024'%20xmlns='http://www.w3.org/2000/svg'%3e%3ctitle%3eFedora%3c/title%3e%3cpath%20d='M12.001%200C5.376%200%20.008%205.369.004%2011.992H.002v9.287h.002A2.726%202.726%200%200%200%202.73%2024h9.275c6.626-.004%2011.993-5.372%2011.993-11.997C23.998%205.375%2018.628%200%2012%200zm2.431%204.94c2.015%200%203.917%201.543%203.917%203.671%200%20.197.001.395-.03.619a1.002%201.002%200%200%201-1.137.893%201.002%201.002%200%200%201-.842-1.175%202.61%202.61%200%200%200%20.013-.337c0-1.207-.987-1.672-1.92-1.672-.934%200-1.775.784-1.777%201.672.016%201.027%200%202.046%200%203.07l1.732-.012c1.352-.028%201.368%202.009.016%201.998l-1.748.013c-.004.826.006.677.002%201.093%200%200%20.015%201.01-.016%201.776-.209%202.25-2.124%204.046-4.424%204.046-2.438%200-4.448-1.993-4.448-4.437.073-2.515%202.078-4.492%204.603-4.469l1.409-.01v1.996l-1.409.013h-.007c-1.388.04-2.577.984-2.6%202.47a2.438%202.438%200%200%200%202.452%202.439c1.356%200%202.441-.987%202.441-2.437l-.001-7.557c0-.14.005-.252.02-.407.23-1.848%201.883-3.256%203.754-3.256z'/%3e%3c/svg%3e`,Ae=`data:image/svg+xml,%3csvg%20role='img'%20viewBox='0%200%2024%2024'%20xmlns='http://www.w3.org/2000/svg'%3e%3ctitle%3eHomebrew%3c/title%3e%3cpath%20d='M7.938%200a.214.214%200%200%200-.206.156c-.316%201.104.179%202.15.838%202.935.153.181.313.347.476.501a2.039%202.039%200%200%200-.665.02c-1.184.233-2.193.985-2.74%202.532a3.893%203.893%200%200%200-.2%201.466%201.565%201.565%200%200%200-1.156%201.504%201.59%201.59%200%200%200%201.227%201.541l.026%2012.046c0%20.195.1.377.264.482a.214.214%200%200%200%20.008.005c.537.31%202.047.812%205.21.812%203.238%200%204.7-.678%205.181-1.04a.214.214%200%200%200%20.008-.007.571.571%200%200%200%20.206-.439c.002-.344.002-1.136.002-1.604a.143.143%200%200%201%20.147-.144c.397.006.869.006%201.318.005a1.826%201.826%200%200%200%201.832-1.825v-5.804a1.826%201.826%200%200%200-1.825-1.826H16.56a.14.14%200%200%201-.143-.144V10.6h.007v-.001a1.573%201.573%200%200%200%201.356-1.556c0-.816-.627-1.489-1.424-1.563-.025-1.438-.437-2.126-.736-2.58a.214.214%200%200%200-.005-.007c-.364-.51-1.193-1.282-2.275-1.316-.503-.016-.842.124-1.125.254-.217.1-.42.177-.67.22.002-1.286.945-1.981.945-1.981a.214.214%200%200%200%20.05-.298s-.087-.122-.21-.26c-.121-.136-.269-.294-.47-.378a.214.214%200%200%200-.079-.017.214.214%200%200%200-.145.055%204.308%204.308%200%200%200-.875%201.101%203.42%203.42%200%200%200-.133.273%203.497%203.497%200%200%200-.381-.846C9.794.978%209.063.436%208.017.016A.214.214%200%200%200%207.939%200zm.156.524c.85.378%201.43.83%201.79%201.403.274.438.426.962.484%201.584a3.07%203.07%200%200%200-.012.462%206.897%206.897%200%200%201-.168-.052%205.487%205.487%200%200%201-1.29-1.106c-.551-.657-.935-1.46-.804-2.291zM11.8%201.618c.07.054.141.101.212.18.034.039.032.04.058.073-.332.308-1.07%201.144-.952%202.453a.214.214%200%200%200%20.222.195c.469-.017.782-.172%201.056-.299.273-.126.508-.228.931-.214.875.027%201.639.715%201.939%201.134.295.449.65%201%20.663%202.36a1.66%201.66%200%200%200-.41.142%201.938%201.938%200%200%200-1.77-1.16%201.94%201.94%200%200%200-1.87%201.448%201.783%201.783%200%200%200-1.356-.64c-.484%200-.91.205-1.233.517a1.873%201.873%200%200%200-1.85-1.625c-.649%200-1.218.335-1.552.84a3.1%203.1%200%200%201%20.157-.735c.51-1.437%201.355-2.045%202.42-2.254.367-.073.664-.011.99.095.325.106.671.262%201.094.342a.214.214%200%200%200%20.252-.245c-.112-.67.073-1.266.336-1.744a3.71%203.71%200%200%201%20.663-.863zM7.44%206.611a1.442%201.442%200%200%201%201.363%201.925.214.214%200%200%200%20.168.283h.005a.214.214%200%200%200%20.238-.146%201.373%201.373%200%200%201%202.613-.01.214.214%200%200%200%20.417-.09%201.509%201.509%200%200%201%201.504-1.664c.678%200%201.249.445%201.442%201.056a.214.214%200%200%200%20.259.143l.15-.04a.214.214%200%200%200%20.051-.02%201.139%201.139%200%200%201%201.702.995%201.14%201.14%200%200%201-.985%201.131.214.214%200%200%200-.001%200%202.215%202.215%200%200%200-.485.126%2010.65%2010.65%200%200%201-1.176.365.214.214%200%200%200-.162.186%201.276%201.276%200%200%201-.146.478%202.07%202.07%200%200%200-.239%201.111l.001.151a.438.438%200%200%201-.16.36.665.665%200%200%201-.43.14.586.586%200%200%201-.588-.59.803.803%200%200%200-.38-.681.214.214%200%200%200-.002-.002c-.24-.145-.43-.37-.532-.636a.214.214%200%200%200-.207-.138%2019.469%2019.469%200%200%201-5.37-.6l-.003-.002a9.007%209.007%200%200%200-.838-.194h.003a1.16%201.16%200%200%201-.937-1.134c0-.619.488-1.118%201.101-1.14a.214.214%200%200%200%20.204-.176%201.443%201.443%200%200%201%201.42-1.187zm8.549%204.106v.455c0%20.314.259.573.572.573h1.329a1.397%201.397%200%200%201%201.397%201.397v5.804a1.396%201.396%200%200%201-1.402%201.396.214.214%200%200%200-.002%200c-.448.002-.918%200-1.31-.005a.573.573%200%200%200-.584.573c0%20.468%200%201.262-.002%201.603a.214.214%200%200%200%200%20.001c0%20.042-.019.08-.05.107-.346.26-1.75.95-4.915.95-3.107%200-4.587-.52-4.99-.752a.143.143%200%200%201-.065-.118l-.025-11.955c.145.033.288.07.431.11a.214.214%200%200%200%20.003%200c.115.031.246.064.383.097v10.37c0%20.129.069.247.18.31.453.217%201.767.732%204.071.732%202.32%200%203.595-.626%204.022-.884a.357.357%200%200%200%20.164-.3l.001-10.21c.267-.075.531-.158.792-.254zm-7.99.894a.493.493%200%200%201%20.494.493v8.578a.493.493%200%200%201-.493.493.493.493%200%200%201-.494-.493v-8.578A.493.493%200%200%201%208%2011.611zm8.652%201.14a.663.663%200%200%200-.662.662v5.208a.663.663%200%200%200%20.662.662h1.14a.663.663%200%200%200%20.662-.662v-5.209a.663.663%200%200%200-.662-.662zm0%20.428h1.14a.233.233%200%200%201%20.233.233v5.21a.233.233%200%200%201-.233.232h-1.14a.233.233%200%200%201-.233-.233v-5.209a.233.233%200%200%201%20.233-.233z'/%3e%3c/svg%3e`,je=`/assets/data-linux.svg`,Me=`data:image/svg+xml,%3csvg%20role='img'%20viewBox='0%200%2024%2024'%20xmlns='http://www.w3.org/2000/svg'%3e%3ctitle%3eManjaro%3c/title%3e%3cpath%20d='M2.182%200A2.177%202.177%200%200%200%200%202.182v19.636C0%2023.027.973%2024%202.182%2024h4.363V6.545h8.728V0Zm15.273%200v24h4.363A2.177%202.177%200%200%200%2024%2021.818V2.182A2.177%202.177%200%200%200%2021.818%200ZM8.727%208.727V24h6.546V8.727Z'/%3e%3c/svg%3e`,Ne=`data:image/svg+xml,%3csvg%20role='img'%20viewBox='0%200%2024%2024'%20xmlns='http://www.w3.org/2000/svg'%3e%3ctitle%3enpm%3c/title%3e%3cpath%20d='M1.763%200C.786%200%200%20.786%200%201.763v20.474C0%2023.214.786%2024%201.763%2024h20.474c.977%200%201.763-.786%201.763-1.763V1.763C24%20.786%2023.214%200%2022.237%200zM5.13%205.323l13.837.019-.009%2013.836h-3.464l.01-10.382h-3.456L12.04%2019.17H5.113z'/%3e%3c/svg%3e`,Pe=`data:image/svg+xml,%3csvg%20role='img'%20viewBox='0%200%2024%2024'%20xmlns='http://www.w3.org/2000/svg'%3e%3ctitle%3epnpm%3c/title%3e%3cpath%20d='M0%200v7.5h7.5V0zm8.25%200v7.5h7.498V0zm8.25%200v7.5H24V0zM2%202h3.5v3.5H2zm8.25%200h3.498v3.5H10.25zm8.25%200H22v3.5h-3.5zM8.25%208.25v7.5h7.498v-7.5zm8.25%200v7.5H24v-7.5zm2%202H22v3.5h-3.5zM0%2016.5V24h7.5v-7.5zm8.25%200V24h7.498v-7.5zm8.25%200V24H24v-7.5z'/%3e%3c/svg%3e`,Fe=`data:image/svg+xml,%3csvg%20role='img'%20viewBox='0%200%2024%2024'%20xmlns='http://www.w3.org/2000/svg'%3e%3ctitle%3eRed%20Hat%3c/title%3e%3cpath%20d='M16.009%2013.386c1.577%200%203.86-.326%203.86-2.202a1.765%201.765%200%200%200-.04-.431l-.94-4.08c-.216-.898-.406-1.305-1.982-2.093-1.223-.625-3.888-1.658-4.676-1.658-.733%200-.947.946-1.822.946-.842%200-1.467-.706-2.255-.706-.757%200-1.25.515-1.63%201.576%200%200-1.06%202.99-1.197%203.424a.81.81%200%200%200-.028.245c0%201.162%204.577%204.974%2010.71%204.974m4.101-1.435c.218%201.032.218%201.14.218%201.277%200%201.765-1.984%202.745-4.593%202.745-5.895.004-11.06-3.451-11.06-5.734a2.326%202.326%200%200%201%20.19-.925C2.746%209.415%200%209.794%200%2012.217c0%203.969%209.405%208.861%2016.851%208.861%205.71%200%207.149-2.582%207.149-4.62%200-1.605-1.387-3.425-3.887-4.512'/%3e%3c/svg%3e`,Ie=`data:image/svg+xml,%3csvg%20role='img'%20viewBox='0%200%2024%2024'%20xmlns='http://www.w3.org/2000/svg'%3e%3ctitle%3eUbuntu%3c/title%3e%3cpath%20d='M17.61.455a3.41%203.41%200%200%200-3.41%203.41%203.41%203.41%200%200%200%203.41%203.41%203.41%203.41%200%200%200%203.41-3.41%203.41%203.41%200%200%200-3.41-3.41zM12.92.8C8.923.777%205.137%202.941%203.148%206.451a4.5%204.5%200%200%201%20.26-.007%204.92%204.92%200%200%201%202.585.737A8.316%208.316%200%200%201%2012.688%203.6%204.944%204.944%200%200%201%2013.723.834%2011.008%2011.008%200%200%200%2012.92.8zm9.226%204.994a4.915%204.915%200%200%201-1.918%202.246%208.36%208.36%200%200%201-.273%208.303%204.89%204.89%200%200%201%201.632%202.54%2011.156%2011.156%200%200%200%20.559-13.089zM3.41%207.932A3.41%203.41%200%200%200%200%2011.342a3.41%203.41%200%200%200%203.41%203.409%203.41%203.41%200%200%200%203.41-3.41%203.41%203.41%200%200%200-3.41-3.41zm2.027%207.866a4.908%204.908%200%200%201-2.915.358%2011.1%2011.1%200%200%200%207.991%206.698%2011.234%2011.234%200%200%200%202.422.249%204.879%204.879%200%200%201-.999-2.85%208.484%208.484%200%200%201-.836-.136%208.304%208.304%200%200%201-5.663-4.32zm11.405.928a3.41%203.41%200%200%200-3.41%203.41%203.41%203.41%200%200%200%203.41%203.41%203.41%203.41%200%200%200%203.41-3.41%203.41%203.41%200%200%200-3.41-3.41z'/%3e%3c/svg%3e`,Le=`data:image/svg+xml,%3csvg%20role='img'%20viewBox='0%200%2024%2024'%20xmlns='http://www.w3.org/2000/svg'%3e%3ctitle%3eYarn%3c/title%3e%3cpath%20d='M12%200C5.375%200%200%205.375%200%2012s5.375%2012%2012%2012%2012-5.375%2012-12S18.625%200%2012%200zm.768%204.105c.183%200%20.363.053.525.157.125.083.287.185.755%201.154.31-.088.468-.042.551-.019.204.056.366.19.463.375.477.917.542%202.553.334%203.605-.241%201.232-.755%202.029-1.131%202.576.324.329.778.899%201.117%201.825.278.774.31%201.478.273%202.015a5.51%205.51%200%200%200%20.602-.329c.593-.366%201.487-.917%202.553-.931.714-.009%201.269.445%201.353%201.103a1.23%201.23%200%200%201-.945%201.362c-.649.158-.95.278-1.821.843-1.232.797-2.539%201.242-3.012%201.39a1.686%201.686%200%200%201-.704.343c-.737.181-3.266.315-3.466.315h-.046c-.783%200-1.214-.241-1.45-.491-.658.329-1.51.19-2.122-.134a1.078%201.078%200%200%201-.58-1.153%201.243%201.243%200%200%201-.153-.195c-.162-.25-.528-.936-.454-1.946.056-.723.556-1.367.88-1.71a5.522%205.522%200%200%201%20.408-2.256c.306-.727.885-1.348%201.32-1.737-.32-.537-.644-1.367-.329-2.21.227-.602.412-.936.82-1.08h-.005c.199-.074.389-.153.486-.259a3.418%203.418%200%200%201%202.298-1.103c.037-.093.079-.185.125-.283.31-.658.639-1.029%201.024-1.168a.94.94%200%200%201%20.328-.06zm.006.7c-.507.016-1.001%201.519-1.001%201.519s-1.27-.204-2.266.871c-.199.218-.468.334-.746.44-.079.028-.176.023-.417.672-.371.991.625%202.094.625%202.094s-1.186.839-1.626%201.881c-.486%201.144-.338%202.261-.338%202.261s-.843.732-.899%201.487c-.051.663.139%201.2.343%201.515.227.343.51.176.51.176s-.561.653-.037.931c.477.25%201.283.394%201.71-.037.31-.31.371-1.001.486-1.283.028-.065.12.111.209.199.097.093.264.195.264.195s-.755.324-.445%201.066c.102.246.468.403%201.066.398.222-.005%202.664-.139%203.313-.296.375-.088.505-.283.505-.283s1.566-.431%202.998-1.357c.917-.598%201.293-.76%202.034-.936.612-.148.57-1.098-.241-1.084-.839.009-1.575.44-2.196.825-1.163.718-1.742.672-1.742.672l-.018-.032c-.079-.13.371-1.293-.134-2.678-.547-1.515-1.413-1.881-1.344-1.997.297-.5%201.038-1.297%201.334-2.78.176-.899.13-2.377-.269-3.151-.074-.144-.732.241-.732.241s-.616-1.371-.788-1.483a.271.271%200%200%200-.157-.046z'/%3e%3c/svg%3e`,Re=`data:image/svg+xml,%3c?xml%20version='1.0'%20encoding='utf-8'?%3e%3c!--%20Generator:%20Adobe%20Illustrator%2021.0.0,%20SVG%20Export%20Plug-In%20.%20SVG%20Version:%206.00%20Build%200)%20--%3e%3csvg%20version='1.1'%20id='Layer_1'%20xmlns='http://www.w3.org/2000/svg'%20xmlns:xlink='http://www.w3.org/1999/xlink'%20x='0px'%20y='0px'%20width='128px'%20height='128px'%20viewBox='0%200%20128%20128'%20enable-background='new%200%200%20128%20128'%20xml:space='preserve'%3e%3cg%3e%3cline%20fill='none'%20x1='0'%20y1='128'%20x2='0'%20y2='0'/%3e%3c/g%3e%3cline%20opacity='0'%20fill-rule='evenodd'%20clip-rule='evenodd'%20fill='%2300FF18'%20x1='0'%20y1='128'%20x2='0'%20y2='0'/%3e%3clinearGradient%20id='SVGID_1_'%20gradientUnits='userSpaceOnUse'%20x1='95.2667'%20y1='91.9263'%20x2='26.7'%20y2='30.68'%3e%3cstop%20offset='0'%20style='stop-color:%23A9C8FF'/%3e%3cstop%20offset='1'%20style='stop-color:%23C7E6FF'/%3e%3c/linearGradient%3e%3cpath%20opacity='0.8'%20fill-rule='evenodd'%20clip-rule='evenodd'%20fill='url(%23SVGID_1_)'%20d='M9.033,109c-1.633,0-3.046-0.638-3.978-1.798%20c-0.952-1.185-1.279-2.814-0.896-4.47l17.986-77.911C22.899,21.557,26.062,19,29.349,19h89.623c1.634,0,3.047,0.638,3.978,1.798%20c0.952,1.184,1.279,2.814,0.896,4.47l-17.986,77.911c-0.753,3.264-3.917,5.822-7.203,5.822H9.033z'/%3e%3cg%3e%3cg%3e%3clinearGradient%20id='SVGID_2_'%20gradientUnits='userSpaceOnUse'%20x1='26.5854'%20y1='30.7778'%20x2='93.5854'%20y2='90.2778'%3e%3cstop%20offset='0'%20style='stop-color:%232D4664'/%3e%3cstop%20offset='0.1689'%20style='stop-color:%2329405B'/%3e%3cstop%20offset='0.4445'%20style='stop-color:%231E2F43'/%3e%3cstop%20offset='0.7902'%20style='stop-color:%230C131B'/%3e%3cstop%20offset='1'%20style='stop-color:%23000000'/%3e%3c/linearGradient%3e%3cpath%20fill-rule='evenodd'%20clip-rule='evenodd'%20fill='url(%23SVGID_2_)'%20d='M118.5,20H29.634c-2.769,0-5.53,2.259-6.168,5.045%20L5.632,102.955C4.995,105.742,6.722,108,9.491,108h88.865c2.769,0,5.53-2.258,6.168-5.045l17.834-77.911%20C122.996,22.259,121.268,20,118.5,20z'/%3e%3c/g%3e%3c/g%3e%3cpath%20fill-rule='evenodd'%20clip-rule='evenodd'%20fill='%232C5591'%20d='M64.165,87.558h21.613c2.513,0,4.55,2.125,4.55,4.746%20c0,2.621-2.037,4.747-4.55,4.747H64.165c-2.513,0-4.55-2.125-4.55-4.747C59.615,89.683,61.652,87.558,64.165,87.558z'/%3e%3cpath%20fill-rule='evenodd'%20clip-rule='evenodd'%20fill='%232C5591'%20d='M78.184,66.455c-0.372,0.749-1.144,1.575-2.509,2.534%20L35.562,97.798c-2.19,1.591-5.334,1.001-7.021-1.319c-1.687-2.32-1.28-5.49,0.91-7.082l36.173-26.194v-0.538L42.896,38.487%20c-1.854-1.972-1.661-5.161,0.431-7.124c2.092-1.962,5.29-1.954,7.144,0.018l27.271,29.012C79.29,62.04,79.405,64.534,78.184,66.455z%20'/%3e%3cpath%20fill-rule='evenodd'%20clip-rule='evenodd'%20fill='%23FFFFFF'%20d='M77.184,65.455c-0.372,0.749-1.144,1.575-2.509,2.534%20L34.562,96.798c-2.19,1.591-5.334,1.001-7.021-1.319c-1.687-2.32-1.28-5.49,0.91-7.082l36.173-26.194v-0.538L41.896,37.487%20c-1.854-1.972-1.661-5.161,0.431-7.124c2.092-1.962,5.29-1.954,7.144,0.018l27.271,29.012C78.29,61.04,78.405,63.534,77.184,65.455z%20'/%3e%3cpath%20fill-rule='evenodd'%20clip-rule='evenodd'%20fill='%23FFFFFF'%20d='M63.55,87h21.613c2.513,0,4.55,2.015,4.55,4.5%20c0,2.485-2.037,4.5-4.55,4.5H63.55C61.037,96,59,93.985,59,91.5C59,89.015,61.037,87,63.55,87z'/%3e%3c/svg%3e`,ze=`/assets/data-scoop.svg`,Be=[`alpine-linux`,`arch-linux`,`curl`,`debian`,`fedora`,`homebrew`,`linux`,`manjaro`,`npm`,`npx`,`pnpm`,`powershell`,`red-hat`,`scoop`,`ubuntu`,`wget`,`yarn`],Ve={"alpine-linux":{kind:`mask`,source:Te,title:`Alpine Linux`},"arch-linux":{kind:`mask`,source:Ee,title:`Arch Linux`},curl:{kind:`mask`,source:De,title:`curl`},debian:{kind:`mask`,source:Oe,title:`Debian`},fedora:{kind:`mask`,source:ke,title:`Fedora`},homebrew:{kind:`mask`,source:Ae,title:`Homebrew`},linux:{kind:`mask`,source:je,title:`Linux`},manjaro:{kind:`mask`,source:Me,title:`Manjaro`},npm:{kind:`mask`,source:Ne,title:`npm`},npx:{kind:`mask`,source:Ne,title:`npx`},pnpm:{kind:`mask`,source:Pe,title:`pnpm`},powershell:{kind:`mask`,source:Re,title:`PowerShell`},"red-hat":{kind:`mask`,source:Fe,title:`Red Hat`},scoop:{kind:`mask`,source:ze,title:`Scoop`},ubuntu:{kind:`mask`,source:Ie,title:`Ubuntu`},wget:{kind:`glyph`,title:`GNU Wget`},yarn:{kind:`mask`,source:Le,title:`Yarn`}};function He(e){return Be.includes(e)?e:`linux`}var Ue=class extends U{constructor(...e){super(...e),this.name=`linux`,this.label=``}static{this.properties={name:{type:String},label:{type:String}}}static{this.styles=o`
    :host { display: inline-flex; width: 1.5rem; height: 1.5rem; color: currentColor; }
    svg, .mask { display: block; width: 100%; height: 100%; }
    .mask { background: currentColor; mask: var(--brand-mask) center / contain no-repeat; -webkit-mask: var(--brand-mask) center / contain no-repeat; }
  `}render(){let e=Ve[He(this.name)],t=this.label||void 0;return e.kind===`glyph`?M`<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="square" stroke-linejoin="miter" role=${t?`img`:`presentation`} aria-label=${t} aria-hidden=${t?`false`:`true`}><path d="M12 3v12m-4-4 4 4 4-4M4 17v4h16v-4"></path></svg>`:M`<span class="mask" style=${`--brand-mask: url("${e.source}")`} role=${t?`img`:`presentation`} aria-label=${t} aria-hidden=${t?`false`:`true`}></span>`}},We=`/assets/data-vector-mark.svg`;function W(e){return Number(e.toFixed(2)).toString()}function Ge(e,t){if(![e.x,e.y,t.x,t.y].every(Number.isFinite))throw TypeError(`Bezier coordinates must be finite`);let n=Math.abs(t.x-e.x),r=Math.max(32,n*.42),i=t.x>=e.x?1:-1,a=e.x+r*i,o=t.x-r*i;return`M ${W(e.x)} ${W(e.y)} C ${W(a)} ${W(e.y)}, ${W(o)} ${W(t.y)}, ${W(t.x)} ${W(t.y)}`}var Ke=class extends U{static{this.styles=o`
    :host { display: inline-flex; min-width: 0; color: var(--courier-color-text, #0e0f0d); font-family: var(--courier-font-sans, sans-serif); }
    .lockup { display: inline-flex; min-width: 0; align-items: center; gap: 0.58rem; color: inherit; }
    img { flex: 0 0 auto; width: 2.2rem; height: 2.2rem; object-fit: contain; }
    strong { font-family: var(--courier-font-display, sans-serif); font-size: 1.08rem; font-weight: 820; letter-spacing: -0.045em; }
  `}render(){return M`<span class="lockup"><img part="mark" src=${We} alt="" decoding="async"><strong part="wordmark">Courier</strong></span>`}},qe=class extends U{constructor(...e){super(...e),this.alt=``,this.eager=!1,this.mobileSource=``,this.source=``}static{this.properties={alt:{type:String},eager:{type:Boolean},mobileSource:{type:String,attribute:`mobile-source`},source:{type:String}}}static{this.styles=o`:host { display: block; } picture { display: contents; } img { display: block; width: 100%; height: auto; }`}render(){return M`<picture>${this.mobileSource?M`<source media="(max-width: 44rem)" srcset=${this.mobileSource}>`:``}<img part="image" src=${this.source} alt=${this.alt} decoding="async" loading=${this.eager?`eager`:`lazy`} fetchpriority=${this.eager?`high`:`auto`}></picture>`}},Je=class extends U{constructor(...e){super(...e),this.tone=`neutral`}static{this.properties={tone:{type:String,reflect:!0}}}static{this.styles=o`
    :host { display: inline-flex; width: max-content; align-items: center; gap: 0.45rem; color: var(--courier-color-muted, #68645e); font-family: var(--courier-font-mono, monospace); font-size: 0.6875rem; font-weight: 700; }
    i { width: 0.48rem; height: 0.48rem; border: 1px solid currentColor; border-radius: 50%; background: currentColor; }
    :host([tone="signal"]) { color: var(--courier-success, #39765d); }
    :host([tone="warning"]) { color: var(--courier-warning, #a86618); }
    :host([tone="danger"]) { color: var(--courier-danger, #b7443e); }
  `}render(){return M`<i aria-hidden="true"></i><slot></slot>`}},Ye=class extends U{constructor(...e){super(...e),this.source=``,this.destination=``}static{this.properties={source:{type:String},destination:{type:String}}}static{this.styles=o`
    :host { display: grid; color: var(--courier-color-text, #0e0f0d); font-family: var(--courier-font-mono, monospace); }
    .route { display: grid; grid-template-columns: minmax(0, 1fr) minmax(3rem, 0.55fr) minmax(0, 1fr); align-items: center; gap: 0.65rem; }
    .node { overflow: hidden; padding: 0.62rem 0; border-bottom: 1px solid var(--courier-color-border, #cbc5b8); font-size: 0.75rem; text-overflow: ellipsis; white-space: nowrap; }
    .node:last-child { text-align: right; }
    .connector { position: relative; height: 2rem; }
    svg { position: absolute; inset: 0; width: 100%; height: 100%; overflow: visible; color: var(--courier-color-border-strong, #8f8a81); }
    path { fill: none; stroke: currentColor; stroke-width: 1.5; vector-effect: non-scaling-stroke; }
    .terminal { position: absolute; top: 50%; width: 0.65rem; height: 0.65rem; aspect-ratio: 1; border: 2px solid var(--courier-color-accent, #ad431d); border-radius: 50%; background: var(--courier-color-canvas, #f2efe6); transform: translateY(-50%); }
    .terminal.source { left: 0; transform: translate(-50%, -50%); }
    .terminal.destination { right: 0; transform: translate(50%, -50%); }
  `}render(){let e=Ge({x:1,y:15},{x:99,y:15});return M`<div class="route"><span class="node">${this.source}</span><span class="connector" aria-hidden="true"><svg viewBox="0 0 100 30" preserveAspectRatio="none"><path d=${e}></path></svg><i class="terminal source"></i><i class="terminal destination"></i></span><span class="node">${this.destination}</span></div>`}},Xe=class extends U{constructor(...e){super(...e),this.checked=!1,this.disabled=!1,this.label=``}static{this.properties={checked:{type:Boolean,reflect:!0},disabled:{type:Boolean,reflect:!0},label:{type:String}}}static{this.styles=o`
    :host { display: inline-flex; min-width: 0; color: inherit; font-family: var(--courier-font-sans, sans-serif); }
    label { display: inline-grid; min-width: 0; grid-template-columns: 1.25rem minmax(0, 1fr); align-items: center; gap: 0.55rem; color: inherit; font-size: 0.75rem; font-weight: 720; line-height: 1.25; cursor: pointer; }
    input { position: absolute; width: 1px; height: 1px; margin: -1px; overflow: hidden; clip-path: inset(50%); white-space: nowrap; }
    .box { position: relative; display: grid; width: 1.25rem; height: 1.25rem; place-items: center; border: 1px solid currentColor; border-radius: 0.25rem; background: color-mix(in srgb, currentColor 6%, transparent); transition: color var(--courier-duration, 160ms) var(--courier-ease, ease), background var(--courier-duration, 160ms) var(--courier-ease, ease), transform var(--courier-duration, 160ms) var(--courier-ease, ease); }
    .box::after { content: ""; width: 0.55rem; height: 0.3rem; border-bottom: 2px solid currentColor; border-left: 2px solid currentColor; opacity: 0; transform: translateY(-0.08rem) rotate(-45deg) scale(0.65); transition: opacity var(--courier-duration, 160ms) var(--courier-ease, ease), transform var(--courier-duration, 160ms) var(--courier-ease, ease); }
    input:checked + .box { border-color: var(--courier-signal, #d4ff45); color: var(--courier-graphite-950, #10120f); background: var(--courier-signal, #d4ff45); }
    input:checked + .box::after { opacity: 1; transform: translateY(-0.08rem) rotate(-45deg) scale(1); }
    input:focus-visible + .box { outline: 3px solid var(--courier-beak, #ff8758); outline-offset: 2px; }
    label:hover .box { transform: translateY(-1px); }
    :host([disabled]) { opacity: 0.52; }
    :host([disabled]) label { cursor: not-allowed; }
    :host([disabled]) label:hover .box { transform: none; }
  `}change(e){this.checked=e.currentTarget.checked,this.dispatchEvent(new CustomEvent(`courier-checkbox-change`,{detail:this.checked,bubbles:!0,composed:!0}))}render(){return M`<label>
      <input type="checkbox" .checked=${this.checked} .disabled=${this.disabled} @change=${this.change}>
      <span class="box" aria-hidden="true"></span>
      <span>${this.label}</span>
    </label>`}},Ze=class extends U{constructor(...e){super(...e),this.href=``,this.icon=`github`,this.label=``,this.target=`_blank`,this.rel=`noopener noreferrer`}static{this.properties={href:{type:String},icon:{type:String},label:{type:String},target:{type:String},rel:{type:String}}}static{this.styles=o`
    :host { display: inline-flex; }
    a {
      box-sizing: border-box;
      display: inline-grid;
      width: var(--courier-control-frame-size, 2.35rem);
      height: var(--courier-control-frame-size, 2.35rem);
      place-items: center;
      border: var(--courier-control-frame-border-width, 1px) solid var(--courier-control-frame-border, var(--courier-color-border, #c8cdbf));
      border-radius: var(--courier-control-frame-radius, var(--courier-radius-md, 0.625rem));
      color: var(--courier-control-frame-color, var(--courier-color-muted, #596054));
      background: var(--courier-control-frame-surface, color-mix(in srgb, var(--courier-color-field, #e7e9dc) 68%, transparent));
      box-shadow: var(--courier-control-frame-shadow, inset 0 1px 2px rgb(16 18 15 / 0.07));
      text-decoration: none;
      transition: color var(--courier-duration, 160ms) var(--courier-ease, ease), border-color var(--courier-duration, 160ms) var(--courier-ease, ease), background var(--courier-duration, 160ms) var(--courier-ease, ease), transform var(--courier-duration, 160ms) var(--courier-ease, ease);
    }
    a:hover { border-color: var(--courier-control-frame-hover-border, var(--courier-color-accent, #d4ff45)); color: var(--courier-control-frame-hover-color, var(--courier-color-text, #151714)); background: var(--courier-control-frame-hover-surface, color-mix(in srgb, var(--courier-color-surface-raised, #fff) 72%, transparent)); }
    a:active { transform: translateY(1px); }
    a:focus-visible { outline: var(--courier-control-frame-focus-width, 3px) solid var(--courier-control-frame-focus, var(--courier-beak, #ff8758)); outline-offset: 2px; }
    courier-icon { width: 1.05rem; height: 1.05rem; }
  `}render(){return M`<a href=${this.href} target=${this.target} rel=${this.rel} aria-label=${this.label} title=${this.label}><courier-icon name=${this.icon}></courier-icon></a>`}},Qe={en:{"theme.label":`Theme`,"theme.system":`System`,"theme.light":`Light`,"theme.dark":`Dark`,"locale.label":`Language`,"locale.en":`English`,"locale.ru":`Russian`,"preference.next":`Next`,"progress.label":`Delivery progress`,"action.cancel":`Cancel`,"action.close":`Close`},ru:{"theme.label":`Тема`,"theme.system":`Системная`,"theme.light":`Светлая`,"theme.dark":`Тёмная`,"locale.label":`Язык`,"locale.en":`Английский`,"locale.ru":`Русский`,"preference.next":`Следующая`,"progress.label":`Ход доставки`,"action.cancel":`Отмена`,"action.close":`Закрыть`}},$e=Object.keys(Qe),et=`courier.locale`;function G(e){if(!e)return;let t=e.toLowerCase().split(`-`)[0];return $e.includes(t)?t:void 0}function tt(e){for(let t of e){let e=G(t);if(e)return e}return`en`}function nt(e,t){if(e)try{let t=G(e.getItem(et));if(t)return t}catch{}return tt(t)}function rt(e,t){if(e)try{e.setItem(et,t)}catch{}}function K(e,t){return Qe[G(e)??`en`][t]}function it(){let e;try{e=globalThis.localStorage}catch{e=void 0}let t=globalThis.navigator?.languages??[];return nt(e,t)}var at=[`system`,`light`,`dark`],ot=`courier.theme`;function st(e){return at.includes(e)?e:`system`}function ct(e,t){return e===`system`?t?.matches?`dark`:`light`:e}function lt(e){if(!e)return`system`;try{return st(e.getItem(ot))}catch{return`system`}}function ut(e,t){if(e)try{e.setItem(ot,t)}catch{}}var dt=class{constructor(e,t,n,r){this.root=e,this.storage=t,this.media=n,this.onSystemChange=()=>this.apply(),this.preference=r??lt(t),this.media?.addEventListener(`change`,this.onSystemChange),this.apply()}set(e){this.preference=e,ut(this.storage,e),this.apply()}destroy(){this.media?.removeEventListener(`change`,this.onSystemChange)}apply(){this.root.dataset.courierTheme=ct(this.preference,this.media),this.root.dataset.courierThemePreference=this.preference}},q=[`system`,`light`,`dark`],J=[`en`,`ru`];function Y(e,t){return e[(e.indexOf(t)+1)%e.length]}var ft=class extends EventTarget{constructor(e,t,n,r,i){super(),this.storage=t,this.themeState=new dt(e,t,n,i),this.locale=nt(t,r)}get theme(){return this.themeState.preference}setTheme(e){return this.themeState.set(e),this.dispatchEvent(new Event(`change`)),e}cycleTheme(){return this.setTheme(Y(q,this.theme))}setLocale(e){return this.locale=e,rt(this.storage,e),this.dispatchEvent(new Event(`change`)),e}cycleLocale(){return this.setLocale(Y(J,this.locale))}destroy(){this.themeState.destroy()}},X;function pt(){if(X)return X;let e;try{e=globalThis.localStorage}catch{e=void 0}let t=globalThis.matchMedia?.(`(prefers-color-scheme: dark)`),n=globalThis.navigator?.languages??[];return X=new ft(document.documentElement,e,t,n),X}var mt={en:`🇬🇧`,ru:`🇷🇺`},ht=class extends U{constructor(...e){super(...e),this.locale=`en`,this.synchronize=()=>{this.controller&&(this.locale=this.controller.locale)}}static{this.properties={locale:{type:String}}}static{this.styles=o`
    :host { display: inline-flex; }
    button { appearance: none; display: inline-grid; width: var(--courier-control-frame-size, 2.35rem); height: var(--courier-control-frame-size, 2.35rem); place-items: center; padding: 0; border: 1px solid var(--courier-control-frame-border, #8f8a81); border-radius: var(--courier-control-frame-radius, 999px); color: var(--courier-control-frame-color, #57544f); background: var(--courier-control-frame-surface, transparent); cursor: pointer; transition: border-color var(--courier-duration, 150ms) var(--courier-ease, ease), background var(--courier-duration, 150ms) var(--courier-ease, ease); }
    button:hover { border-color: var(--courier-control-frame-hover-border, #ad431d); background: var(--courier-control-frame-hover-surface, transparent); }
    button:focus-visible { outline: var(--courier-control-frame-focus-width, 3px) solid var(--courier-control-frame-focus, #d95f2b); outline-offset: 2px; }
    span { font: 1.02rem/1 system-ui, sans-serif; }
  `}connectedCallback(){super.connectedCallback(),this.controller=pt(),this.synchronize(),this.controller.addEventListener(`change`,this.synchronize)}disconnectedCallback(){this.controller?.removeEventListener(`change`,this.synchronize),super.disconnectedCallback()}cycle(){let e=this.controller?.cycleLocale()??Y(J,this.locale);this.locale=e,this.dispatchEvent(new CustomEvent(`courier-locale-change`,{detail:e,bubbles:!0,composed:!0}))}render(){let e=Y(J,this.locale),t=`${K(this.locale,`locale.label`)}: ${K(this.locale,`locale.${this.locale}`)}. ${K(this.locale,`preference.next`)}: ${K(this.locale,`locale.${e}`)}`;return M`<button type="button" aria-label=${t} title=${t} @click=${this.cycle}><span aria-hidden="true">${mt[this.locale]}</span></button>`}},gt=class extends U{constructor(...e){super(...e),this.heading=``}static{this.properties={heading:{type:String}}}static{this.styles=o`
    :host {
      display: block;
      overflow: hidden;
      border: 1px solid var(--courier-color-border, #c8cdbf);
      border-radius: var(--courier-radius-md, 0.5rem);
      color: var(--courier-color-text, #151714);
      background: var(--courier-color-surface-raised, #fff);
      font-family: var(--courier-font-sans, sans-serif);
      box-shadow: 0 1px 0 rgb(16 18 15 / 0.04);
    }
    section { padding: var(--courier-space-6, 1.5rem); }
    h2 {
      margin: 0 0 var(--courier-space-4, 1rem);
      font-family: var(--courier-font-mono, monospace);
      font-size: 0.6875rem;
      letter-spacing: 0.09em;
      text-transform: uppercase;
    }
  `}render(){return this.heading?M`<section aria-labelledby="courier-panel-heading">
      <h2 id="courier-panel-heading">${this.heading}</h2>
      <slot></slot>
    </section>`:M`<section><slot></slot></section>`}};function _t(e,t){return!Number.isFinite(e)||!Number.isFinite(t)||t<=0?0:Math.min(1,Math.max(0,e/t))}var vt=class extends U{constructor(...e){super(...e),this.value=0,this.total=0,this.label=``,this.locale=`en`}static{this.properties={value:{type:Number},total:{type:Number},label:{type:String},locale:{type:String}}}static{this.styles=o`
    :host {
      display: grid;
      gap: var(--courier-space-2, 0.5rem);
      color: var(--courier-color-text, #151714);
      font-family: var(--courier-font-sans, sans-serif);
    }
    .track {
      overflow: hidden;
      height: 0.5rem;
      border: 1px solid var(--courier-color-border, #c8cdbf);
      border-radius: var(--courier-radius-xs, 0.125rem);
      background: var(--courier-color-field, #e7e9dc);
    }
    .fill {
      height: 100%;
      background: var(--courier-color-accent, #d4ff45);
      transform-origin: left;
      transition: transform var(--courier-duration, 160ms) var(--courier-ease, ease);
    }
    output { color: var(--courier-color-muted, #596054); font-family: var(--courier-font-mono, monospace); font-size: 0.75rem; }
  `}render(){let e=_t(this.value,this.total),t=this.label||K(this.locale,`progress.label`),n=Number.isFinite(this.total)&&this.total>0?this.total:0;return M`<div
      class="track"
      role="progressbar"
      aria-label=${t}
      aria-valuemin="0"
      aria-valuemax=${n}
      aria-valuenow=${Number.isFinite(this.value)?Math.max(0,Math.min(this.value,n)):0}
    ><div class="fill" style=${`transform: scaleX(${e})`}></div></div>
    <output>${Math.round(e*100)}%</output>`}},yt=class extends U{constructor(...e){super(...e),this.eager=!1,this.mobileSource=``,this.source=``,this.active=!1}static{this.properties={eager:{type:Boolean},mobileSource:{type:String,attribute:`mobile-source`},source:{type:String},active:{state:!0}}}static{this.styles=o`
    :host { position: absolute; display: block; inset: 0; overflow: hidden; background: var(--courier-color-canvas, #f2efe6); pointer-events: none; isolation: isolate; }
    courier-mascot { position: absolute; inset: 0; width: 100%; height: 100%; pointer-events: none; }
    courier-mascot::part(image) { width: 100%; height: 100%; object-fit: cover; }
    .veil { position: absolute; inset: 0; background: linear-gradient(180deg, rgb(14 15 13 / 0.02), transparent 26% 76%, rgb(14 15 13 / 0.08)); pointer-events: none; }
  `}connectedCallback(){super.connectedCallback(),this.configureLoading()}disconnectedCallback(){this.proximityObserver?.disconnect(),this.proximityObserver=void 0,super.disconnectedCallback()}updated(e){e.has(`eager`)&&this.eager&&this.activate()}configureLoading(){if(!this.active){if(this.eager){this.activate();return}if(!globalThis.IntersectionObserver){this.active=!0;return}this.proximityObserver=new IntersectionObserver(e=>{e.some(e=>e.isIntersecting)&&this.activate()},{rootMargin:`100% 0px`,threshold:0}),this.proximityObserver.observe(this)}}activate(){this.active=!0,this.proximityObserver?.disconnect(),this.proximityObserver=void 0}render(){return M`${this.active?M`<courier-mascot class="base" ?eager=${this.eager} alt="" .source=${this.source} .mobileSource=${this.mobileSource}></courier-mascot>`:``}<span class="veil" aria-hidden="true"></span>`}},bt={system:`system`,light:`sun`,dark:`moon`},xt=class extends U{constructor(...e){super(...e),this.preference=`system`,this.locale=`en`,this.synchronize=()=>{this.controller&&(this.preference=this.controller.theme)}}static{this.properties={preference:{type:String},locale:{type:String}}}static{this.styles=o`
    :host { display: inline-flex; }
    button { appearance: none; display: inline-grid; width: var(--courier-control-frame-size, 2.35rem); height: var(--courier-control-frame-size, 2.35rem); place-items: center; padding: 0; border: 1px solid var(--courier-control-frame-border, #8f8a81); border-radius: var(--courier-control-frame-radius, 999px); color: var(--courier-control-frame-color, #57544f); background: var(--courier-control-frame-surface, transparent); cursor: pointer; transition: color var(--courier-duration, 150ms) var(--courier-ease, ease), border-color var(--courier-duration, 150ms) var(--courier-ease, ease), background var(--courier-duration, 150ms) var(--courier-ease, ease); }
    button:hover { color: var(--courier-control-frame-hover-color, #0e0f0d); border-color: var(--courier-control-frame-hover-border, #ad431d); background: var(--courier-control-frame-hover-surface, transparent); }
    button:focus-visible { outline: var(--courier-control-frame-focus-width, 3px) solid var(--courier-control-frame-focus, #d95f2b); outline-offset: 2px; }
    courier-icon { width: 1.05rem; height: 1.05rem; }
  `}connectedCallback(){super.connectedCallback(),this.controller=pt(),this.synchronize(),this.controller.addEventListener(`change`,this.synchronize)}disconnectedCallback(){this.controller?.removeEventListener(`change`,this.synchronize),super.disconnectedCallback()}cycle(){let e=this.controller?.cycleTheme()??Y(q,this.preference);this.preference=e,this.dispatchEvent(new CustomEvent(`courier-theme-change`,{detail:e,bubbles:!0,composed:!0}))}render(){let e=Y(q,this.preference),t=K(this.locale,`theme.${this.preference}`),n=K(this.locale,`theme.${e}`),r=`${K(this.locale,`theme.label`)}: ${t}. ${K(this.locale,`preference.next`)}: ${n}`;return M`<button type="button" aria-label=${r} title=${r} @click=${this.cycle}><courier-icon name=${bt[this.preference]}></courier-icon></button>`}},St=class extends U{constructor(...e){super(...e),this.heading=``,this.status=``}static{this.properties={heading:{type:String},status:{type:String}}}static{this.styles=o`
    :host { display: grid; min-width: 0; min-height: 0; grid-template-rows: auto minmax(0, 1fr) auto; overflow: hidden; border-top: 1px solid var(--courier-workbench-line, #cbc5b8); border-bottom: 1px solid var(--courier-workbench-line, #cbc5b8); color: var(--courier-color-text, #0e0f0d); background: var(--courier-workbench-surface, rgb(242 239 230 / 0.78)); backdrop-filter: blur(16px); }
    header { display: flex; min-height: 2.75rem; align-items: center; justify-content: space-between; gap: 1rem; padding: 0.55rem 0.2rem; border-bottom: 1px solid var(--courier-workbench-line, #cbc5b8); color: var(--courier-color-muted, #68645e); font-size: 0.72rem; font-weight: 720; }
    .status { color: var(--courier-color-accent, #ad431d); font-family: var(--courier-font-mono, monospace); }
    .body { min-width: 0; min-height: 0; overflow: auto; }
    footer { border-top: 1px solid var(--courier-workbench-line, #cbc5b8); }
  `}render(){return M`<header><span>${this.heading}</span><slot name="actions"></slot>${this.status?M`<span class="status">${this.status}</span>`:P}</header><div class="body"><slot></slot></div><footer><slot name="footer"></slot></footer>`}},Ct=class extends U{constructor(...e){super(...e),this.command=``,this.description=``,this.heading=`Command`,this.copyLabel=`Copy`,this.copiedLabel=`Copied`,this.copyFailedLabel=`Copy failed`,this.sessionKey=``,this.copyState=`idle`}static{this.properties={command:{type:String},description:{type:String},heading:{type:String},copyLabel:{type:String,attribute:`copy-label`},copiedLabel:{type:String,attribute:`copied-label`},copyFailedLabel:{type:String,attribute:`copy-failed-label`},sessionKey:{type:String,attribute:`session-key`},copyState:{state:!0}}}static{this.styles=o`
    :host { display: grid; min-width: 0; min-height: 0; grid-template-rows: auto minmax(0, 1fr) auto; border-top: 1px solid var(--courier-workbench-line, #cbc5b8); border-bottom: 1px solid var(--courier-workbench-line, #cbc5b8); color: var(--courier-color-text, #0e0f0d); background: var(--courier-workbench-surface, rgb(242 239 230 / 0.78)); backdrop-filter: blur(16px); }
    header { display: flex; min-height: 2.75rem; align-items: center; justify-content: space-between; gap: 1rem; padding: 0.48rem 0.2rem; border-bottom: 1px solid var(--courier-workbench-line, #cbc5b8); color: var(--courier-color-muted, #68645e); font-size: 0.72rem; font-weight: 720; }
    button { appearance: none; display: inline-flex; min-height: 2rem; align-items: center; gap: 0.38rem; padding: 0.3rem 0.65rem; border: 1px solid var(--courier-control-frame-border, #8f8a81); border-radius: 999px; color: var(--courier-color-text, #0e0f0d); background: transparent; font: 700 0.68rem/1 var(--courier-font-sans, sans-serif); cursor: pointer; }
    button:hover:not(:disabled) { border-color: var(--courier-color-accent, #ad431d); color: var(--courier-color-accent, #ad431d); }
    button:focus-visible { outline: 3px solid var(--courier-color-accent, #ad431d); outline-offset: 2px; }
    button:disabled { cursor: not-allowed; opacity: 0.42; }
    courier-icon { width: 0.9rem; height: 0.9rem; }
    .readout { display: grid; min-width: 0; min-height: 0; grid-template-rows: auto auto minmax(0, 1fr); }
    .command, .description { min-width: 0; padding: 0.75rem 0.2rem; }
    .command { color: var(--courier-color-text, #0e0f0d); }
    code { overflow-wrap: anywhere; font: 700 0.78rem/1.5 var(--courier-font-mono, monospace); white-space: pre-wrap; }
    .description { margin: 0; border-top: 1px solid color-mix(in srgb, var(--courier-workbench-line, #cbc5b8) 65%, transparent); color: var(--courier-color-muted, #68645e); font-size: 0.72rem; line-height: 1.45; }
    .details { min-width: 0; min-height: 0; overflow: auto; }
    footer { display: flex; min-width: 0; min-height: 2.4rem; align-items: center; justify-content: space-between; gap: 0.75rem; padding: 0.4rem 0.2rem; border-top: 1px solid var(--courier-workbench-line, #cbc5b8); color: var(--courier-color-muted, #68645e); font-size: 0.62rem; }
    .copy-status { min-width: 6rem; min-height: 1.4em; }
    @media (max-width: 44rem) { header button span { display: none; } .command, .description { padding: 0.6rem 0.1rem; } footer { padding-inline: 0.1rem; } }
  `}willUpdate(e){(e.has(`sessionKey`)||e.has(`command`))&&this.reset()}reset(){this.copyState=`idle`}async copyCommand(){let e=this.clipboard??globalThis.navigator.clipboard;if(!e||!this.command){this.copyState=`failed`;return}try{await e.writeText(this.command),this.copyState=`copied`}catch{this.copyState=`failed`}}render(){let e=this.copyState===`copied`?this.copiedLabel:this.copyState===`failed`?this.copyFailedLabel:``;return M`<header><span>${this.heading}</span><button type="button" ?disabled=${!this.command} aria-label=${this.copyLabel} title=${this.copyLabel} @click=${this.copyCommand}><courier-icon name=${this.copyState===`copied`?`check`:`copy`}></courier-icon><span>${this.copyLabel}</span></button></header><div class="readout"><div class="command"><code>${this.command}</code></div>${this.description?M`<p class="description">${this.description}</p>`:P}<div class="details"><slot name="details"></slot></div></div><footer><span class="copy-status" role="status" aria-live="polite">${e}</span><slot name="footer-actions"></slot></footer>`}},wt=`apple.archive.browser-share.browser-upload.check.copy.download.folder.folder-in.folder-out.github.homebrew.linux.moon.npm.package.parcel.pnpm.receipt.retry.route.scoop.server.server-in.server-out.shield.sun.system.terminal.upload.webhook-in.webhook-out.windows.yarn`.split(`.`),Tt={apple:`M15 5c1-1 1-3 1-3-2 0-3 1-4 3m6 7c-1-2-2-3-4-3-1 0-2 1-3 1s-2-1-3-1c-3 0-5 3-5 6 0 4 3 8 5 8 1 0 2-1 3-1s2 1 3 1c2 0 4-3 5-6-2-1-3-2-3-4 0-2 1-3 2-4z`,archive:`M3 3h18v5H3zM5 8v13h14V8M9 12h6`,"browser-share":`M3 4h18v15H3zM3 8h18M7 6h.01M10 6h.01M14 15c2-3 4-4 7-4m-3-2 3 2-1 4`,"browser-upload":`M3 4h18v15H3zM3 8h18M7 6h.01M10 6h.01M12 17v-6m-3 3 3-3 3 3`,check:`m4 12 5 5L20 6`,copy:`M8 3h13v13M3 8h13v13H3z`,download:`M12 3v13m-5-5 5 5 5-5M4 15v6h16v-6`,folder:`M3 6h7l2 3h9v12H3zM3 6V3h7l2 3h9v3`,"folder-in":`M3 6h7l2 3h9v12H3zM3 6V3h7l2 3h9v3m12 6h-7m3-3-3 3 3 3`,"folder-out":`M3 6h7l2 3h9v12H3zM3 6V3h7l2 3h9v3m5 6h7m-3-3 3 3-3 3`,github:`M9 19c-5 1-5-2-7-3m14 6v-3.6c0-1 .1-1.7-.4-2.2 3.2-.4 6.4-1.6 6.4-7.1 0-1.6-.6-3-1.7-4 .2-.5.7-2.3-.2-4.6 0 0-1.4-.5-4.7 1.7a16 16 0 0 0-8.6 0C6.4 1 5 1.5 5 1.5 4.1 3.8 4.6 5.6 4.8 6.1a7 7 0 0 0-1.7 4c0 5.5 3.2 6.7 6.4 7.1-.4.4-.8 1.1-.8 2.2V23`,homebrew:`M6 4h11l-1 15H8zM17 7h2a2 2 0 0 1 2 2v3a2 2 0 0 1-2 2h-2M5 22h13`,linux:`M12 2c-3 0-4 3-4 6-2 2-3 5-3 8l3-1 1 5 3-2 3 2 1-5 3 1c0-3-1-6-3-8 0-3-1-6-4-6zM9 8h.01M15 8h.01M10 11h4`,moon:`M20 16a8 8 0 0 1-12-10 8 8 0 1 0 12 10z`,npm:`M2 6h20v12H2zM6 15V9h5v6m0-6h4v6m0-6h3v6`,package:`m3 7 9-5 9 5v11l-9 4-9-4zM3 7l9 5 9-5M12 12v10M8 4l9 5`,parcel:`m3 7 9-5 9 5v11l-9 4-9-4zM3 7l9 5 9-5M12 12v10M8 4l9 5v5`,pnpm:`M3 3h5v5H3zM10 3h5v5h-5zM17 3h4v5h-4zM3 10h5v5H3zm7 0h5v5h-5zm7 0h4v5h-4zM10 17h5v4h-5zm7 0h4v4h-4z`,receipt:`M5 2h14v20l-3-2-4 2-4-2-3 2zM8 7h8M8 11h8m-8 5 2 2 5-4`,retry:`M3 10a9 9 0 1 1 1 7M3 3v7h7M12 7v5l3 2`,route:`M2 3h6v6H2zM16 15h6v6h-6zM11 6h8v6m-3-3 3 3 3-3M13 18H5v-6m-3 3 3-3 3 3`,scoop:`M5 8h14l-2 13H7zM4 8h16M8 8V5a4 4 0 0 1 8 0v3`,server:`M3 2h18v8H3zM3 14h18v8H3zM7 6h1m3 0h6M7 18h1m3 0h6M6 10v4m12-4v4`,"server-in":`M3 2h18v8H3zM3 14h18v8H3zM7 6h1m3 0h6M7 18h1m7-12h-6m3-3-3 3 3 3`,"server-out":`M3 2h18v8H3zM3 14h18v8H3zM7 6h1m3 0h6M7 18h1m3-12h6m-3-3 3 3-3 3`,shield:`m12 2 8 3v7c0 5-8 10-8 10S4 17 4 12V5zM8 11l3 3 5-6`,sun:`M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8zM12 2v3m0 14v3M4.9 4.9 7 7m10 10 2.1 2.1M2 12h3m14 0h3M4.9 19.1 7 17M17 7l2.1-2.1`,system:`M3 4h18v13H3zM8 21h8M12 17v4`,terminal:`M4 6l5 6-5 6m7 0h9`,upload:`M12 16V3m-5 5 5-5 5 5M4 15v6h16v-6`,"webhook-in":`M7 7a3 3 0 1 1 3 3l-3 6a3 3 0 1 0 3 4m7-3a3 3 0 1 1-3-3l3-6a3 3 0 1 0-3-4m-1 8h-6m3-3-3 3 3 3`,"webhook-out":`M7 7a3 3 0 1 1 3 3l-3 6a3 3 0 1 0 3 4m7-3a3 3 0 1 1-3-3l3-6a3 3 0 1 0-3-4m-4 8h6m-3-3 3 3-3 3`,windows:`M3 4l8-1v8H3zm10-1 8-1v9h-8zM3 13h8v8l-8-1zm10 0h8v9l-8-1z`,yarn:`M12 3a9 9 0 1 0 9 9M8 16c4-1 7-4 9-8m-8 1c3 1 5 4 5 8m-5-5c-1-3 0-5 2-6`};function Et(e){return wt.includes(e)?e:`parcel`}var Dt={"courier-brand":Ke,"courier-brand-icon":Ue,"courier-button":we,"courier-checkbox":Xe,"courier-icon":class extends U{constructor(...e){super(...e),this.name=`parcel`,this.label=``}static{this.properties={name:{type:String},label:{type:String}}}static{this.styles=o`
    :host {
      display: inline-flex;
      width: 1.5rem;
      height: 1.5rem;
      color: currentColor;
    }
    svg { width: 100%; height: 100%; }
  `}render(){let e=Et(this.name);return M`<svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="1.7"
      stroke-linecap="square"
      stroke-linejoin="miter"
      role=${this.label?`img`:`presentation`}
      aria-hidden=${this.label?`false`:`true`}
      aria-label=${this.label||void 0}
    ><path d=${Tt[e]}></path></svg>`}},"courier-icon-link":Ze,"courier-locale-selector":ht,"courier-mascot":qe,"courier-panel":gt,"courier-progress":vt,"courier-route":Ye,"courier-scene":yt,"courier-status":Je,"courier-theme-selector":xt,"courier-workbench":St,"courier-command-readout":Ct};function Ot(e=customElements){for(let[t,n]of Object.entries(Dt))e.get(t)||e.define(t,n)}var kt=`/assets/data-delivery-access-mobile-v1.webp`,At=`/assets/data-delivery-access-wide-v1.webp`,jt=kt;function Z(e,t=``,n=globalThis.location.pathname){let r=n.endsWith(`/`)?n:`${n}/`,i=new URL(`api/v1/${e}`,globalThis.location.origin);return i.pathname=`${r}api/v1/${e}`,t&&i.searchParams.set(`path`,t),`${i.pathname}${i.search}`}async function Q(e){if(!e.ok)throw Error(`Courier request failed (${e.status})`);return e.json()}async function Mt(e=``,t=globalThis.fetch){return Q(await t(Z(`meta`,e),{credentials:`same-origin`}))}async function Nt(e,t=globalThis.fetch){return Q(await t(Z(`session`),{method:`POST`,credentials:`same-origin`,headers:{"Content-Type":`application/json`},body:JSON.stringify({password:e})}))}async function Pt(e,t,n=globalThis.fetch){let r=new FormData;r.append(`file`,e),await Q(await n(Z(`upload`),{method:`POST`,credentials:`same-origin`,headers:{"X-Courier-CSRF":t,"X-Courier-File-Size":String(e.size)},body:r}))}function $(e,t=!1){let n=new URL(Z(`download`,e),globalThis.location.origin);return t&&n.searchParams.set(`archive`,`tar.gz`),`${n.pathname}${n.search}`}function Ft(e,t){return e?`${e}/${t}`:t}function It(e){let t=e.lastIndexOf(`/`);return t<0?``:e.slice(0,t)}var Lt={en:{title:`Courier delivery`,privateRoute:`Private route`,loading:`Preparing the delivery route…`,retry:`Retry connection`,download:`Download file`,downloadArchive:`Download as archive`,downloadAll:`Download directory`,upload:`Choose a file`,uploadTitle:`Dispatch a file`,uploadHelp:`Select one file. Courier checks policy and reserves the final path before committing it.`,up:`Parent directory`,password:`Delivery password`,signIn:`Verify access`,accessTitle:`Identity check`,accessHelp:`This route is protected. Verify access to reveal its delivery metadata.`,empty:`No entries are available at this path.`,failed:`The route is unavailable or authorization is required. No delivery metadata was revealed.`,ready:`Route ready`,manifest:`Delivery manifest`,confirmed:`Verified handoff`,entryTypeFile:`File`,entryTypeDirectory:`Directory`,itemSize:`Bytes`},ru:{title:`Доставка Courier`,privateRoute:`Приватный маршрут`,loading:`Подготовка маршрута доставки…`,retry:`Повторить подключение`,download:`Скачать файл`,downloadArchive:`Скачать архивом`,downloadAll:`Скачать каталог`,upload:`Выбрать файл`,uploadTitle:`Отправить файл`,uploadHelp:`Выберите один файл. Courier проверит правила и зарезервирует конечный путь до фиксации.`,up:`Родительский каталог`,password:`Пароль доставки`,signIn:`Подтвердить доступ`,accessTitle:`Проверка доступа`,accessHelp:`Маршрут защищён. Подтвердите доступ, чтобы увидеть данные доставки.`,empty:`По этому пути нет доступных объектов.`,failed:`Маршрут недоступен или требуется авторизация. Данные доставки не были раскрыты.`,ready:`Маршрут готов`,manifest:`Манифест доставки`,confirmed:`Подтверждённая передача`,entryTypeFile:`Файл`,entryTypeDirectory:`Каталог`,itemSize:`Байт`}};function Rt(e,t){return Lt[e][t]}Ot();var zt=class extends U{constructor(...e){super(...e),this.locale=it(),this.failed=!1,this.csrf=``}static{this.properties={locale:{state:!0},metadata:{state:!0},failed:{state:!0},csrf:{state:!0}}}static{this.styles=[Ce,o`
    :host {
      display: block;
      min-height: 100vh;
      padding: 0 1rem 3rem;
      color: var(--courier-color-text);
      background-color: var(--courier-color-canvas);
      background-image: linear-gradient(var(--courier-color-grid) 1px, transparent 1px), linear-gradient(90deg, var(--courier-color-grid) 1px, transparent 1px);
      background-size: 2.5rem 2.5rem;
      font-family: var(--courier-font-sans);
    }
    main, courier-panel { min-width: 0; }
    main { width: min(66rem, 100%); margin: 0 auto; }
    header { display: flex; min-height: 5rem; align-items: center; justify-content: space-between; gap: 1rem; border-bottom: 1px solid var(--courier-color-border); }
    nav, .controls, .actions, .row { display: flex; align-items: center; gap: 0.75rem; flex-wrap: wrap; }
    .workspace { display: grid; gap: 1rem; padding-top: clamp(2rem, 6vw, 5rem); }
    .operation-head { display: flex; align-items: end; justify-content: space-between; gap: 1rem; padding-bottom: 1rem; border-bottom: 1px solid var(--courier-color-border); }
    .operation-head > div { display: grid; gap: 0.45rem; }
    .eyebrow, .label { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; font-weight: 750; letter-spacing: 0.08em; text-transform: uppercase; }
    h1, h2, p { margin: 0; overflow-wrap: anywhere; }
    h1 { font-family: var(--courier-font-display); font-size: clamp(2.4rem, 6vw, 4.75rem); font-weight: 830; letter-spacing: -0.06em; line-height: 0.95; }
    h2 { font-size: clamp(1.35rem, 4vw, 2rem); letter-spacing: -0.035em; }
    p { line-height: 1.6; }
    .muted { color: var(--courier-color-muted); }
    .access { display: grid; grid-template-columns: minmax(0, 1fr) minmax(11rem, 0.45fr); gap: 1rem; overflow: hidden; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-lg); background: var(--courier-color-surface-raised); box-shadow: var(--courier-shadow); }
    .access-copy { display: grid; align-content: center; gap: 1rem; padding: clamp(1.5rem, 5vw, 3.5rem); }
    .access-art { position: relative; min-height: 24rem; overflow: hidden; background: var(--courier-graphite-900); }
    .access-art::before { content: ""; position: absolute; inset: 0; opacity: 0.15; background-image: linear-gradient(rgb(243 244 233 / 0.2) 1px, transparent 1px), linear-gradient(90deg, rgb(243 244 233 / 0.2) 1px, transparent 1px); background-size: 2rem 2rem; }
    .access-art courier-mascot { position: absolute; right: -15%; bottom: -2%; width: 130%; }
    .error { padding: 0.85rem 1rem; border-left: 3px solid var(--courier-warning); color: var(--courier-color-text); background: color-mix(in srgb, var(--courier-warning) 12%, transparent); }
    form { display: grid; gap: 0.75rem; }
    .signin { grid-template-columns: minmax(0, 1fr) auto; }
    .field { display: grid; min-width: 0; gap: 0.35rem; color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; font-weight: 700; letter-spacing: 0.06em; text-transform: uppercase; }
    button.link:focus-visible, a:focus-visible { outline: 3px solid var(--courier-beak); outline-offset: 2px; }
    .route-overview { display: grid; gap: 0.75rem; padding: 1rem; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-md); background: var(--courier-color-surface); }
    .delivery-panel { display: grid; gap: 1.25rem; padding: clamp(1.25rem, 4vw, 2rem); border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-lg); background: var(--courier-color-surface-raised); box-shadow: var(--courier-shadow); }
    .delivery-title { display: flex; align-items: start; justify-content: space-between; gap: 1rem; }
    .delivery-title > div { display: grid; gap: 0.35rem; }
    .toolbar { display: flex; gap: 0.75rem; align-items: center; flex-wrap: wrap; padding: 0.9rem 0; border-top: 1px solid var(--courier-color-border); border-bottom: 1px solid var(--courier-color-border); }
    a, button.link { color: var(--courier-color-text); font-weight: 750; }
    button.link { appearance: none; padding: 0; border: 0; background: transparent; font: inherit; text-decoration: underline; cursor: pointer; }
    .upload-zone { display: grid; gap: 0.75rem; padding: clamp(1.25rem, 4vw, 2rem); border: 1px dashed var(--courier-color-border-strong); border-radius: var(--courier-radius-md); background: var(--courier-color-surface); }
    .upload-zone .courier-file-action { justify-self: start; }
    ul { margin: 0; padding: 0; list-style: none; border-top: 1px solid var(--courier-color-border); }
    li { display: grid; grid-template-columns: auto minmax(0, 1fr) auto auto; gap: 0.8rem; align-items: center; min-height: 3.75rem; padding: 0.75rem 0; border-bottom: 1px solid var(--courier-color-border); }
    li courier-icon { color: var(--courier-color-muted); }
    .entry-name { min-width: 0; overflow-wrap: anywhere; }
    .size { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.75rem; font-variant-numeric: tabular-nums; }
    .loading { display: grid; min-height: 14rem; place-items: center; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-lg); color: var(--courier-color-muted); background: var(--courier-color-surface-raised); font-family: var(--courier-font-mono); }
    .page-scene { position: fixed; z-index: 0; inset: 0; opacity: 0.38; }
    main { position: relative; z-index: 2; }
    :host::after { content: ""; position: fixed; z-index: 1; inset: 0; background: linear-gradient(90deg, var(--courier-color-canvas) 0 24%, color-mix(in srgb, var(--courier-color-canvas) 72%, transparent) 62%, color-mix(in srgb, var(--courier-color-canvas) 88%, transparent)); pointer-events: none; }
    header { border-bottom-color: color-mix(in srgb, var(--courier-color-border) 65%, transparent); }
    .operation-head { border-bottom: 0; }
    .access { display: block; min-height: 0; }
    .access-copy { min-height: 20rem; color: var(--courier-color-text); }
    .access-art { display: none; }
    .route-overview { padding: 0.85rem; border: 0; border-radius: 0; background: transparent; }
    .delivery-panel { display: block; padding: 0; border: 0; border-radius: 0; background: transparent; box-shadow: none; }
    .delivery-title, .toolbar, .upload-zone, .delivery-panel > ul, .delivery-panel > p { margin: 0.75rem; }
    .delivery-title { color: var(--courier-color-text); }
    .toolbar { border-color: var(--courier-workbench-line); }
    .upload-zone { border-color: var(--courier-workbench-line); color: var(--courier-color-text); background: color-mix(in srgb, var(--courier-color-surface) 64%, transparent); }
    ul { border-top-color: var(--courier-workbench-line); }
    li { border-bottom-color: var(--courier-workbench-line); }
    a, button.link { color: var(--courier-color-text); }
    .muted, .size { color: var(--courier-color-muted); }
    .loading { min-height: 10rem; padding: 1rem; }
    @media (max-width: 44rem) {
      header { align-items: flex-start; padding: 1rem 0; }
      nav { justify-content: flex-end; }
      .access { grid-template-columns: 1fr; }
      .access-art { min-height: 18rem; }
      .access-art courier-mascot { right: -4%; bottom: -8%; width: 106%; }
      .signin { grid-template-columns: 1fr; }
      li { grid-template-columns: auto minmax(0, 1fr) auto; }
      li .size { display: none; }
      .delivery-title { display: grid; }
    }
  `]}connectedCallback(){super.connectedCallback(),this.refresh()}disconnectedCallback(){super.disconnectedCallback()}async refresh(e=this.metadata?.path??``){this.failed=!1;try{this.metadata=await Mt(e)}catch{this.failed=!0,this.metadata=void 0}}async openDirectory(e,t){e.preventDefault(),await this.refresh(t)}async signIn(e){e.preventDefault();let t=e.currentTarget,n=new FormData(t).get(`password`)?.toString()??``;try{this.csrf=(await Nt(n)).csrf,t.reset(),await this.refresh()}catch{this.failed=!0}}async sendFile(e){let t=e.currentTarget,n=t.files?.item(0);if(n)try{await Pt(n,this.csrf),t.value=``,await this.refresh()}catch{this.failed=!0}}setLocale(e){this.locale=e.detail}t(e){return Rt(this.locale,e)}entry(e){let t=Ft(this.metadata.path,e.name);return e.type===`directory`?M`<li><courier-icon name="folder"></courier-icon><button class="link entry-name" @click=${e=>this.openDirectory(e,t)}>${e.name}</button><span class="size">${e.size} ${this.t(`itemSize`)}</span><a href=${$(t,!0)}>${this.t(`downloadArchive`)}</a></li>`:M`<li><courier-icon name="parcel"></courier-icon><span class="entry-name">${e.name}</span><span class="size">${e.size} ${this.t(`itemSize`)}</span><a href=${$(t)}>${this.t(`download`)}</a></li>`}render(){let e=this.metadata?.entries??[];return M`
      <courier-scene class="page-scene" eager .source=${At} .mobileSource=${jt}></courier-scene>
      <main>
        <header>
          <courier-brand></courier-brand>
          <nav><courier-theme-selector .locale=${this.locale}></courier-theme-selector><courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector></nav>
        </header>
        <div class="workspace">
          <div class="operation-head"><div><span class="eyebrow">${this.t(`privateRoute`)}</span><h1>${this.t(`title`)}</h1></div>${this.metadata?M`<courier-status tone="signal">${this.t(`ready`)}</courier-status>`:P}</div>
          ${this.failed?M`
            <courier-workbench class="access" .heading=${this.t(`privateRoute`)} status="authentication required">
              <div class="access-copy">
                <span class="eyebrow">${this.t(`privateRoute`)}</span>
                <h2>${this.t(`accessTitle`)}</h2>
                <p class="muted">${this.t(`accessHelp`)}</p>
                <p class="error" role="alert">${this.t(`failed`)}</p>
                <form class="signin" @submit=${this.signIn}><label class="field"><span>${this.t(`password`)}</span><input name="password" type="password" autocomplete="current-password" placeholder=${this.t(`password`)}></label><courier-button type="submit" variant="primary">${this.t(`signIn`)}</courier-button></form>
                <courier-button @click=${this.refresh}>${this.t(`retry`)}</courier-button>
              </div>
            </courier-workbench>
          `:P}
          ${this.metadata?M`
            <courier-workbench class="route-overview" .heading=${this.t(`confirmed`)} status="verified"><courier-route source="sender" destination=${this.metadata.name}></courier-route></courier-workbench>
            <courier-workbench class="delivery-panel" .heading=${this.t(`manifest`)} .status=${this.t(`ready`)}>
              <div class="delivery-title"><div><span class="eyebrow">${this.t(`manifest`)}</span><h2>${this.metadata.name}</h2></div><courier-status tone="signal">${this.t(`ready`)}</courier-status></div>
              ${this.metadata.type===`upload`?M`
                <div class="upload-zone"><h2>${this.t(`uploadTitle`)}</h2><p class="muted">${this.t(`uploadHelp`)}</p><label class="courier-file-action"><courier-icon name="upload"></courier-icon><span>${this.t(`upload`)}</span><input type="file" @change=${this.sendFile}></label></div>
              `:M`
                ${this.metadata.type===`file`?M`<div class="toolbar"><courier-icon name="download"></courier-icon><a href=${$(this.metadata.path)}>${this.t(`download`)}</a></div>`:M`
                  <div class="toolbar"><a href=${$(this.metadata.path,!0)}>${this.t(`downloadAll`)}</a>${this.metadata.path?M`<button class="link" @click=${e=>this.openDirectory(e,It(this.metadata.path))}>${this.t(`up`)}</button>`:P}</div>
                  ${e.length===0?M`<p class="muted">${this.t(`empty`)}</p>`:M`<ul>${e.map(e=>this.entry(e))}</ul>`}
                `}
              `}
            </courier-workbench>
          `:this.failed?P:M`<courier-workbench class="loading" .heading=${this.t(`privateRoute`)} status="running"><courier-status>${this.t(`loading`)}</courier-status></courier-workbench>`}
        </div>
      </main>
    `}};customElements.get(`courier-data-app`)||customElements.define(`courier-data-app`,zt);