(function(){let e=document.createElement(`link`).relList;if(e&&e.supports&&e.supports(`modulepreload`))return;for(let e of document.querySelectorAll(`link[rel="modulepreload"]`))n(e);new MutationObserver(e=>{for(let t of e)if(t.type===`childList`)for(let e of t.addedNodes)e.tagName===`LINK`&&e.rel===`modulepreload`&&n(e)}).observe(document,{childList:!0,subtree:!0});function t(e){let t={};return e.integrity&&(t.integrity=e.integrity),e.referrerPolicy&&(t.referrerPolicy=e.referrerPolicy),t.credentials=e.crossOrigin===`use-credentials`?`include`:e.crossOrigin===`anonymous`?`omit`:`same-origin`,t}function n(e){if(e.ep)return;e.ep=!0;let n=t(e);fetch(e.href,n)}})();var e=globalThis,t=e.ShadowRoot&&(e.ShadyCSS===void 0||e.ShadyCSS.nativeShadow)&&`adoptedStyleSheets`in Document.prototype&&`replace`in CSSStyleSheet.prototype,n=Symbol(),r=new WeakMap,i=class{constructor(e,t,r){if(this._$cssResult$=!0,r!==n)throw Error("CSSResult is not constructable. Use `unsafeCSS` or `css` instead.");this.cssText=e,this.t=t}get styleSheet(){let e=this.o,n=this.t;if(t&&e===void 0){let t=n!==void 0&&n.length===1;t&&(e=r.get(n)),e===void 0&&((this.o=e=new CSSStyleSheet).replaceSync(this.cssText),t&&r.set(n,e))}return e}toString(){return this.cssText}},a=e=>new i(typeof e==`string`?e:e+``,void 0,n),o=(e,...t)=>new i(e.length===1?e[0]:t.reduce((t,n,r)=>t+(e=>{if(!0===e._$cssResult$)return e.cssText;if(typeof e==`number`)return e;throw Error(`Value passed to 'css' function must be a 'css' function result: `+e+`. Use 'unsafeCSS' to pass non-literal values, but take care to ensure page security.`)})(n)+e[r+1],e[0]),e,n),s=(n,r)=>{if(t)n.adoptedStyleSheets=r.map(e=>e instanceof CSSStyleSheet?e:e.styleSheet);else for(let t of r){let r=document.createElement(`style`),i=e.litNonce;i!==void 0&&r.setAttribute(`nonce`,i),r.textContent=t.cssText,n.appendChild(r)}},c=t?e=>e:e=>e instanceof CSSStyleSheet?(e=>{let t=``;for(let n of e.cssRules)t+=n.cssText;return a(t)})(e):e,{is:l,defineProperty:u,getOwnPropertyDescriptor:d,getOwnPropertyNames:ee,getOwnPropertySymbols:te,getPrototypeOf:ne}=Object,f=globalThis,p=f.trustedTypes,re=p?p.emptyScript:``,ie=f.reactiveElementPolyfillSupport,m=(e,t)=>e,h={toAttribute(e,t){switch(t){case Boolean:e=e?re:null;break;case Object:case Array:e=e==null?e:JSON.stringify(e)}return e},fromAttribute(e,t){let n=e;switch(t){case Boolean:n=e!==null;break;case Number:n=e===null?null:Number(e);break;case Object:case Array:try{n=JSON.parse(e)}catch{n=null}}return n}},g=(e,t)=>!l(e,t),_={attribute:!0,type:String,converter:h,reflect:!1,useDefault:!1,hasChanged:g};Symbol.metadata??=Symbol(`metadata`),f.litPropertyMetadata??=new WeakMap;var v=class extends HTMLElement{static addInitializer(e){this._$Ei(),(this.l??=[]).push(e)}static get observedAttributes(){return this.finalize(),this._$Eh&&[...this._$Eh.keys()]}static createProperty(e,t=_){if(t.state&&(t.attribute=!1),this._$Ei(),this.prototype.hasOwnProperty(e)&&((t=Object.create(t)).wrapped=!0),this.elementProperties.set(e,t),!t.noAccessor){let n=Symbol(),r=this.getPropertyDescriptor(e,n,t);r!==void 0&&u(this.prototype,e,r)}}static getPropertyDescriptor(e,t,n){let{get:r,set:i}=d(this.prototype,e)??{get(){return this[t]},set(e){this[t]=e}};return{get:r,set(t){let a=r?.call(this);i?.call(this,t),this.requestUpdate(e,a,n)},configurable:!0,enumerable:!0}}static getPropertyOptions(e){return this.elementProperties.get(e)??_}static _$Ei(){if(this.hasOwnProperty(m(`elementProperties`)))return;let e=ne(this);e.finalize(),e.l!==void 0&&(this.l=[...e.l]),this.elementProperties=new Map(e.elementProperties)}static finalize(){if(this.hasOwnProperty(m(`finalized`)))return;if(this.finalized=!0,this._$Ei(),this.hasOwnProperty(m(`properties`))){let e=this.properties,t=[...ee(e),...te(e)];for(let n of t)this.createProperty(n,e[n])}let e=this[Symbol.metadata];if(e!==null){let t=litPropertyMetadata.get(e);if(t!==void 0)for(let[e,n]of t)this.elementProperties.set(e,n)}this._$Eh=new Map;for(let[e,t]of this.elementProperties){let n=this._$Eu(e,t);n!==void 0&&this._$Eh.set(n,e)}this.elementStyles=this.finalizeStyles(this.styles)}static finalizeStyles(e){let t=[];if(Array.isArray(e)){let n=new Set(e.flat(1/0).reverse());for(let e of n)t.unshift(c(e))}else e!==void 0&&t.push(c(e));return t}static _$Eu(e,t){let n=t.attribute;return!1===n?void 0:typeof n==`string`?n:typeof e==`string`?e.toLowerCase():void 0}constructor(){super(),this._$Ep=void 0,this.isUpdatePending=!1,this.hasUpdated=!1,this._$Em=null,this._$Ev()}_$Ev(){this._$ES=new Promise(e=>this.enableUpdating=e),this._$AL=new Map,this._$E_(),this.requestUpdate(),this.constructor.l?.forEach(e=>e(this))}addController(e){(this._$EO??=new Set).add(e),this.renderRoot!==void 0&&this.isConnected&&e.hostConnected?.()}removeController(e){this._$EO?.delete(e)}_$E_(){let e=new Map,t=this.constructor.elementProperties;for(let n of t.keys())this.hasOwnProperty(n)&&(e.set(n,this[n]),delete this[n]);e.size>0&&(this._$Ep=e)}createRenderRoot(){let e=this.shadowRoot??this.attachShadow(this.constructor.shadowRootOptions);return s(e,this.constructor.elementStyles),e}connectedCallback(){this.renderRoot??=this.createRenderRoot(),this.enableUpdating(!0),this._$EO?.forEach(e=>e.hostConnected?.())}enableUpdating(e){}disconnectedCallback(){this._$EO?.forEach(e=>e.hostDisconnected?.())}attributeChangedCallback(e,t,n){this._$AK(e,n)}_$ET(e,t){let n=this.constructor.elementProperties.get(e),r=this.constructor._$Eu(e,n);if(r!==void 0&&!0===n.reflect){let i=(n.converter?.toAttribute===void 0?h:n.converter).toAttribute(t,n.type);this._$Em=e,i==null?this.removeAttribute(r):this.setAttribute(r,i),this._$Em=null}}_$AK(e,t){let n=this.constructor,r=n._$Eh.get(e);if(r!==void 0&&this._$Em!==r){let e=n.getPropertyOptions(r),i=typeof e.converter==`function`?{fromAttribute:e.converter}:e.converter?.fromAttribute===void 0?h:e.converter;this._$Em=r;let a=i.fromAttribute(t,e.type);this[r]=a??this._$Ej?.get(r)??a,this._$Em=null}}requestUpdate(e,t,n,r=!1,i){if(e!==void 0){let a=this.constructor;if(!1===r&&(i=this[e]),n??=a.getPropertyOptions(e),!((n.hasChanged??g)(i,t)||n.useDefault&&n.reflect&&i===this._$Ej?.get(e)&&!this.hasAttribute(a._$Eu(e,n))))return;this.C(e,t,n)}!1===this.isUpdatePending&&(this._$ES=this._$EP())}C(e,t,{useDefault:n,reflect:r,wrapped:i},a){n&&!(this._$Ej??=new Map).has(e)&&(this._$Ej.set(e,a??t??this[e]),!0!==i||a!==void 0)||(this._$AL.has(e)||(this.hasUpdated||n||(t=void 0),this._$AL.set(e,t)),!0===r&&this._$Em!==e&&(this._$Eq??=new Set).add(e))}async _$EP(){this.isUpdatePending=!0;try{await this._$ES}catch(e){Promise.reject(e)}let e=this.scheduleUpdate();return e!=null&&await e,!this.isUpdatePending}scheduleUpdate(){return this.performUpdate()}performUpdate(){if(!this.isUpdatePending)return;if(!this.hasUpdated){if(this.renderRoot??=this.createRenderRoot(),this._$Ep){for(let[e,t]of this._$Ep)this[e]=t;this._$Ep=void 0}let e=this.constructor.elementProperties;if(e.size>0)for(let[t,n]of e){let{wrapped:e}=n,r=this[t];!0!==e||this._$AL.has(t)||r===void 0||this.C(t,void 0,n,r)}}let e=!1,t=this._$AL;try{e=this.shouldUpdate(t),e?(this.willUpdate(t),this._$EO?.forEach(e=>e.hostUpdate?.()),this.update(t)):this._$EM()}catch(t){throw e=!1,this._$EM(),t}e&&this._$AE(t)}willUpdate(e){}_$AE(e){this._$EO?.forEach(e=>e.hostUpdated?.()),this.hasUpdated||(this.hasUpdated=!0,this.firstUpdated(e)),this.updated(e)}_$EM(){this._$AL=new Map,this.isUpdatePending=!1}get updateComplete(){return this.getUpdateComplete()}getUpdateComplete(){return this._$ES}shouldUpdate(e){return!0}update(e){this._$Eq&&=this._$Eq.forEach(e=>this._$ET(e,this[e])),this._$EM()}updated(e){}firstUpdated(e){}};v.elementStyles=[],v.shadowRootOptions={mode:`open`},v[m(`elementProperties`)]=new Map,v[m(`finalized`)]=new Map,ie?.({ReactiveElement:v}),(f.reactiveElementVersions??=[]).push(`2.1.2`);var y=globalThis,b=e=>e,x=y.trustedTypes,S=x?x.createPolicy(`lit-html`,{createHTML:e=>e}):void 0,ae=`$lit$`,C=`lit$${Math.random().toFixed(9).slice(2)}$`,oe=`?`+C,se=`<${oe}>`,w=document,T=()=>w.createComment(``),E=e=>e===null||typeof e!=`object`&&typeof e!=`function`,D=Array.isArray,ce=e=>D(e)||typeof e?.[Symbol.iterator]==`function`,O=`[ 	
\f\r]`,k=/<(?:(!--|\/[^a-zA-Z])|(\/?[a-zA-Z][^>\s]*)|(\/?$))/g,A=/-->/g,j=/>/g,M=RegExp(`>|${O}(?:([^\\s"'>=/]+)(${O}*=${O}*(?:[^ \t\n\f\r"'\`<>=]|("|')|))|$)`,`g`),N=/'/g,P=/"/g,F=/^(?:script|style|textarea|title)$/i,I=(e=>(t,...n)=>({_$litType$:e,strings:t,values:n}))(1),L=Symbol.for(`lit-noChange`),R=Symbol.for(`lit-nothing`),z=new WeakMap,B=w.createTreeWalker(w,129);function V(e,t){if(!D(e)||!e.hasOwnProperty(`raw`))throw Error(`invalid template strings array`);return S===void 0?t:S.createHTML(t)}var le=(e,t)=>{let n=e.length-1,r=[],i,a=t===2?`<svg>`:t===3?`<math>`:``,o=k;for(let t=0;t<n;t++){let n=e[t],s,c,l=-1,u=0;for(;u<n.length&&(o.lastIndex=u,c=o.exec(n),c!==null);)u=o.lastIndex,o===k?c[1]===`!--`?o=A:c[1]===void 0?c[2]===void 0?c[3]!==void 0&&(o=M):(F.test(c[2])&&(i=RegExp(`</`+c[2],`g`)),o=M):o=j:o===M?c[0]===`>`?(o=i??k,l=-1):c[1]===void 0?l=-2:(l=o.lastIndex-c[2].length,s=c[1],o=c[3]===void 0?M:c[3]===`"`?P:N):o===P||o===N?o=M:o===A||o===j?o=k:(o=M,i=void 0);let d=o===M&&e[t+1].startsWith(`/>`)?` `:``;a+=o===k?n+se:l>=0?(r.push(s),n.slice(0,l)+ae+n.slice(l)+C+d):n+C+(l===-2?t:d)}return[V(e,a+(e[n]||`<?>`)+(t===2?`</svg>`:t===3?`</math>`:``)),r]},H=class e{constructor({strings:t,_$litType$:n},r){let i;this.parts=[];let a=0,o=0,s=t.length-1,c=this.parts,[l,u]=le(t,n);if(this.el=e.createElement(l,r),B.currentNode=this.el.content,n===2||n===3){let e=this.el.content.firstChild;e.replaceWith(...e.childNodes)}for(;(i=B.nextNode())!==null&&c.length<s;){if(i.nodeType===1){if(i.hasAttributes())for(let e of i.getAttributeNames())if(e.endsWith(ae)){let t=u[o++],n=i.getAttribute(e).split(C),r=/([.?@])?(.*)/.exec(t);c.push({type:1,index:a,name:r[2],strings:n,ctor:r[1]===`.`?de:r[1]===`?`?fe:r[1]===`@`?pe:G}),i.removeAttribute(e)}else e.startsWith(C)&&(c.push({type:6,index:a}),i.removeAttribute(e));if(F.test(i.tagName)){let e=i.textContent.split(C),t=e.length-1;if(t>0){i.textContent=x?x.emptyScript:``;for(let n=0;n<t;n++)i.append(e[n],T()),B.nextNode(),c.push({type:2,index:++a});i.append(e[t],T())}}}else if(i.nodeType===8){if(i.data===oe)c.push({type:2,index:a});else{let e=-1;for(;(e=i.data.indexOf(C,e+1))!==-1;)c.push({type:7,index:a}),e+=C.length-1}}a++}}static createElement(e,t){let n=w.createElement(`template`);return n.innerHTML=e,n}};function U(e,t,n=e,r){if(t===L)return t;let i=r===void 0?n._$Cl:n._$Co?.[r],a=E(t)?void 0:t._$litDirective$;return i?.constructor!==a&&(i?._$AO?.(!1),a===void 0?i=void 0:(i=new a(e),i._$AT(e,n,r)),r===void 0?n._$Cl=i:(n._$Co??=[])[r]=i),i!==void 0&&(t=U(e,i._$AS(e,t.values),i,r)),t}var ue=class{constructor(e,t){this._$AV=[],this._$AN=void 0,this._$AD=e,this._$AM=t}get parentNode(){return this._$AM.parentNode}get _$AU(){return this._$AM._$AU}u(e){let{el:{content:t},parts:n}=this._$AD,r=(e?.creationScope??w).importNode(t,!0);B.currentNode=r;let i=B.nextNode(),a=0,o=0,s=n[0];for(;s!==void 0;){if(a===s.index){let t;s.type===2?t=new W(i,i.nextSibling,this,e):s.type===1?t=new s.ctor(i,s.name,s.strings,this,e):s.type===6&&(t=new me(i,this,e)),this._$AV.push(t),s=n[++o]}a!==s?.index&&(i=B.nextNode(),a++)}return B.currentNode=w,r}p(e){let t=0;for(let n of this._$AV)n!==void 0&&(n.strings===void 0?n._$AI(e[t]):(n._$AI(e,n,t),t+=n.strings.length-2)),t++}},W=class e{get _$AU(){return this._$AM?._$AU??this._$Cv}constructor(e,t,n,r){this.type=2,this._$AH=R,this._$AN=void 0,this._$AA=e,this._$AB=t,this._$AM=n,this.options=r,this._$Cv=r?.isConnected??!0}get parentNode(){let e=this._$AA.parentNode,t=this._$AM;return t!==void 0&&e?.nodeType===11&&(e=t.parentNode),e}get startNode(){return this._$AA}get endNode(){return this._$AB}_$AI(e,t=this){e=U(this,e,t),E(e)?e===R||e==null||e===``?(this._$AH!==R&&this._$AR(),this._$AH=R):e!==this._$AH&&e!==L&&this._(e):e._$litType$===void 0?e.nodeType===void 0?ce(e)?this.k(e):this._(e):this.T(e):this.$(e)}O(e){return this._$AA.parentNode.insertBefore(e,this._$AB)}T(e){this._$AH!==e&&(this._$AR(),this._$AH=this.O(e))}_(e){this._$AH!==R&&E(this._$AH)?this._$AA.nextSibling.data=e:this.T(w.createTextNode(e)),this._$AH=e}$(e){let{values:t,_$litType$:n}=e,r=typeof n==`number`?this._$AC(e):(n.el===void 0&&(n.el=H.createElement(V(n.h,n.h[0]),this.options)),n);if(this._$AH?._$AD===r)this._$AH.p(t);else{let e=new ue(r,this),n=e.u(this.options);e.p(t),this.T(n),this._$AH=e}}_$AC(e){let t=z.get(e.strings);return t===void 0&&z.set(e.strings,t=new H(e)),t}k(t){D(this._$AH)||(this._$AH=[],this._$AR());let n=this._$AH,r,i=0;for(let a of t)i===n.length?n.push(r=new e(this.O(T()),this.O(T()),this,this.options)):r=n[i],r._$AI(a),i++;i<n.length&&(this._$AR(r&&r._$AB.nextSibling,i),n.length=i)}_$AR(e=this._$AA.nextSibling,t){for(this._$AP?.(!1,!0,t);e!==this._$AB;){let t=b(e).nextSibling;b(e).remove(),e=t}}setConnected(e){this._$AM===void 0&&(this._$Cv=e,this._$AP?.(e))}},G=class{get tagName(){return this.element.tagName}get _$AU(){return this._$AM._$AU}constructor(e,t,n,r,i){this.type=1,this._$AH=R,this._$AN=void 0,this.element=e,this.name=t,this._$AM=r,this.options=i,n.length>2||n[0]!==``||n[1]!==``?(this._$AH=Array(n.length-1).fill(new String),this.strings=n):this._$AH=R}_$AI(e,t=this,n,r){let i=this.strings,a=!1;if(i===void 0)e=U(this,e,t,0),a=!E(e)||e!==this._$AH&&e!==L,a&&(this._$AH=e);else{let r=e,o,s;for(e=i[0],o=0;o<i.length-1;o++)s=U(this,r[n+o],t,o),s===L&&(s=this._$AH[o]),a||=!E(s)||s!==this._$AH[o],s===R?e=R:e!==R&&(e+=(s??``)+i[o+1]),this._$AH[o]=s}a&&!r&&this.j(e)}j(e){e===R?this.element.removeAttribute(this.name):this.element.setAttribute(this.name,e??``)}},de=class extends G{constructor(){super(...arguments),this.type=3}j(e){this.element[this.name]=e===R?void 0:e}},fe=class extends G{constructor(){super(...arguments),this.type=4}j(e){this.element.toggleAttribute(this.name,!!e&&e!==R)}},pe=class extends G{constructor(e,t,n,r,i){super(e,t,n,r,i),this.type=5}_$AI(e,t=this){if((e=U(this,e,t,0)??R)===L)return;let n=this._$AH,r=e===R&&n!==R||e.capture!==n.capture||e.once!==n.once||e.passive!==n.passive,i=e!==R&&(n===R||r);r&&this.element.removeEventListener(this.name,this,n),i&&this.element.addEventListener(this.name,this,e),this._$AH=e}handleEvent(e){typeof this._$AH==`function`?this._$AH.call(this.options?.host??this.element,e):this._$AH.handleEvent(e)}},me=class{constructor(e,t,n){this.element=e,this.type=6,this._$AN=void 0,this._$AM=t,this.options=n}get _$AU(){return this._$AM._$AU}_$AI(e){U(this,e)}},he=y.litHtmlPolyfillSupport;he?.(H,W),(y.litHtmlVersions??=[]).push(`3.3.3`);var ge=(e,t,n)=>{let r=n?.renderBefore??t,i=r._$litPart$;if(i===void 0){let e=n?.renderBefore??null;r._$litPart$=i=new W(t.insertBefore(T(),e),e,void 0,n??{})}return i._$AI(e),i},K=globalThis,q=class extends v{constructor(){super(...arguments),this.renderOptions={host:this},this._$Do=void 0}createRenderRoot(){let e=super.createRenderRoot();return this.renderOptions.renderBefore??=e.firstChild,e}update(e){let t=this.render();this.hasUpdated||(this.renderOptions.isConnected=this.isConnected),super.update(e),this._$Do=ge(t,this.renderRoot,this.renderOptions)}connectedCallback(){super.connectedCallback(),this._$Do?.setConnected(!0)}disconnectedCallback(){super.disconnectedCallback(),this._$Do?.setConnected(!1)}render(){return L}};q._$litElement$=!0,q.finalized=!0,K.litElementHydrateSupport?.({LitElement:q});var _e=K.litElementPolyfillSupport;_e?.({LitElement:q}),(K.litElementVersions??=[]).push(`4.2.2`);var J=o`
  :host {
    color: var(--courier-text, #151714);
    font-family: var(--courier-font-sans, sans-serif);
  }

  button,
  select {
    min-height: 2.75rem;
    border: 1px solid var(--courier-line, #c7ccc0);
    border-radius: var(--courier-radius-sm, 0.25rem);
    color: inherit;
    background: var(--courier-surface, #fff);
    font: inherit;
  }

  button {
    padding: 0.625rem 1rem;
    cursor: pointer;
    transition: background var(--courier-duration, 160ms) var(--courier-ease, ease);
  }

  button:hover:not(:disabled) {
    background: color-mix(in srgb, var(--courier-signal, #d4ff45) 24%, var(--courier-surface, #fff));
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
    outline: 3px solid var(--courier-beak, #ff8758);
    outline-offset: 2px;
  }
`,ve=o`
  label {
    display: grid;
    gap: var(--courier-space-1, 0.25rem);
    color: var(--courier-muted, #51574d);
    font-size: 0.8125rem;
    font-weight: 700;
    letter-spacing: 0.04em;
    text-transform: uppercase;
  }
`,ye=class extends q{constructor(...e){super(...e),this.disabled=!1,this.variant=`secondary`}static{this.properties={disabled:{type:Boolean,reflect:!0},variant:{type:String,reflect:!0}}}static{this.styles=[J,o`
    :host { display: inline-flex; }
    button { width: 100%; }
    :host([variant="primary"]) button {
      border-color: var(--courier-signal, #d4ff45);
      color: var(--courier-signal-ink, #151714);
      background: var(--courier-signal, #d4ff45);
      font-weight: 800;
    }
  `]}render(){return I`<button type="button" ?disabled=${this.disabled}><slot></slot></button>`}},be={en:{"theme.label":`Theme`,"theme.system":`System`,"theme.light":`Light`,"theme.dark":`Dark`,"locale.label":`Language`,"locale.en":`English`,"locale.ru":`Russian`,"progress.label":`Delivery progress`,"action.cancel":`Cancel`,"action.close":`Close`},ru:{"theme.label":`Тема`,"theme.system":`Системная`,"theme.light":`Светлая`,"theme.dark":`Тёмная`,"locale.label":`Язык`,"locale.en":`Английский`,"locale.ru":`Русский`,"progress.label":`Ход доставки`,"action.cancel":`Отмена`,"action.close":`Закрыть`}},xe=Object.keys(be),Y=`courier.locale`;function X(e){if(!e)return;let t=e.toLowerCase().split(`-`)[0];return xe.includes(t)?t:void 0}function Se(e){for(let t of e){let e=X(t);if(e)return e}return`en`}function Ce(e,t){if(e)try{let t=X(e.getItem(Y));if(t)return t}catch{}return Se(t)}function we(e,t){if(e)try{e.setItem(Y,t)}catch{}}function Z(e,t){return be[X(e)??`en`][t]}function Te(){let e;try{e=globalThis.localStorage}catch{e=void 0}let t=globalThis.navigator?.languages??[];return Ce(e,t)}var Ee=class extends q{constructor(...e){super(...e),this.locale=`en`}static{this.properties={locale:{type:String}}}static{this.styles=[J,ve]}connectedCallback(){super.connectedCallback(),this.locale=Te()}change(e){let t=X(e.currentTarget.value)??`en`;this.locale=t;let n;try{n=globalThis.localStorage}catch{n=void 0}we(n,t),this.dispatchEvent(new CustomEvent(`courier-locale-change`,{detail:t,bubbles:!0,composed:!0}))}render(){return I`<label>${Z(this.locale,`locale.label`)}
      <select .value=${this.locale} @change=${this.change} data-storage-key=${Y}>
        <option value="en">${Z(this.locale,`locale.en`)}</option>
        <option value="ru">${Z(this.locale,`locale.ru`)}</option>
      </select>
    </label>`}},De=class extends q{constructor(...e){super(...e),this.heading=``}static{this.properties={heading:{type:String}}}static{this.styles=o`
    :host {
      display: block;
      border: 1px solid var(--courier-line, #c7ccc0);
      border-radius: var(--courier-radius-md, 0.5rem);
      color: var(--courier-text, #151714);
      background: var(--courier-surface, #fff);
      font-family: var(--courier-font-sans, sans-serif);
    }
    section { padding: var(--courier-space-6, 1.5rem); }
    h2 { margin: 0 0 var(--courier-space-4, 1rem); font-size: 1.125rem; }
  `}render(){return I`<section aria-labelledby="courier-panel-heading">
      <h2 id="courier-panel-heading">${this.heading}</h2>
      <slot></slot>
    </section>`}};function Oe(e,t){return!Number.isFinite(e)||!Number.isFinite(t)||t<=0?0:Math.min(1,Math.max(0,e/t))}var ke=class extends q{constructor(...e){super(...e),this.value=0,this.total=0,this.label=``,this.locale=`en`}static{this.properties={value:{type:Number},total:{type:Number},label:{type:String},locale:{type:String}}}static{this.styles=o`
    :host {
      display: grid;
      gap: var(--courier-space-2, 0.5rem);
      color: var(--courier-text, #151714);
      font-family: var(--courier-font-sans, sans-serif);
    }
    .track {
      overflow: hidden;
      height: 0.625rem;
      border: 1px solid var(--courier-line, #c7ccc0);
      border-radius: 999px;
      background: var(--courier-surface, #fff);
    }
    .fill {
      height: 100%;
      background: var(--courier-signal, #d4ff45);
      transform-origin: left;
      transition: transform var(--courier-duration, 160ms) var(--courier-ease, ease);
    }
    output { color: var(--courier-muted, #51574d); font-family: var(--courier-font-mono, monospace); }
  `}render(){let e=Oe(this.value,this.total),t=this.label||Z(this.locale,`progress.label`),n=Number.isFinite(this.total)&&this.total>0?this.total:0;return I`<div
      class="track"
      role="progressbar"
      aria-label=${t}
      aria-valuemin="0"
      aria-valuemax=${n}
      aria-valuenow=${Number.isFinite(this.value)?Math.max(0,Math.min(this.value,n)):0}
    ><div class="fill" style=${`transform: scaleX(${e})`}></div></div>
    <output>${Math.round(e*100)}%</output>`}},Ae=[`system`,`light`,`dark`],je=`courier.theme`;function Me(e){return Ae.includes(e)?e:`system`}function Ne(e,t){return e===`system`?t?.matches?`dark`:`light`:e}function Pe(e){if(!e)return`system`;try{return Me(e.getItem(je))}catch{return`system`}}function Fe(e,t){if(e)try{e.setItem(je,t)}catch{}}var Ie=class{constructor(e,t,n,r){this.root=e,this.storage=t,this.media=n,this.onSystemChange=()=>this.apply(),this.preference=r??Pe(t),this.media?.addEventListener(`change`,this.onSystemChange),this.apply()}set(e){this.preference=e,Fe(this.storage,e),this.apply()}destroy(){this.media?.removeEventListener(`change`,this.onSystemChange)}apply(){this.root.dataset.courierTheme=Ne(this.preference,this.media),this.root.dataset.courierThemePreference=this.preference}};function Le(){let e;try{e=globalThis.localStorage}catch{e=void 0}let t=globalThis.matchMedia?.(`(prefers-color-scheme: dark)`);return new Ie(document.documentElement,e,t)}var Re=class extends q{constructor(...e){super(...e),this.preference=`system`,this.locale=`en`}static{this.properties={preference:{type:String},locale:{type:String}}}static{this.styles=[J,ve]}connectedCallback(){super.connectedCallback(),this.state=Le(),this.preference=this.state.preference}disconnectedCallback(){this.state?.destroy(),super.disconnectedCallback()}change(e){let t=Me(e.currentTarget.value);this.preference=t,this.state?.set(t),this.dispatchEvent(new CustomEvent(`courier-theme-change`,{detail:t,bubbles:!0,composed:!0}))}render(){return I`<label>${Z(this.locale,`theme.label`)}
      <select .value=${this.preference} @change=${this.change}>
        <option value="system">${Z(this.locale,`theme.system`)}</option>
        <option value="light">${Z(this.locale,`theme.light`)}</option>
        <option value="dark">${Z(this.locale,`theme.dark`)}</option>
      </select>
    </label>`}},ze=[`archive`,`copy`,`download`,`folder`,`parcel`,`receipt`,`retry`,`route`,`server`,`shield`,`upload`],Be={archive:`M3 3h18v5H3zM5 8v13h14V8M9 12h6`,copy:`M8 3h13v13M3 8h13v13H3z`,download:`M12 3v13m-5-5 5 5 5-5M4 15v6h16v-6`,folder:`M3 6h7l2 3h9v12H3zM3 6V3h7l2 3h9v3`,parcel:`m3 7 9-5 9 5v11l-9 4-9-4zM3 7l9 5 9-5M12 12v10M8 4l9 5v5`,receipt:`M5 2h14v20l-3-2-4 2-4-2-3 2zM8 7h8M8 11h8m-8 5 2 2 5-4`,retry:`M3 10a9 9 0 1 1 1 7M3 3v7h7M12 7v5l3 2`,route:`M2 3h6v6H2zM16 15h6v6h-6zM11 6h8v6m-3-3 3 3 3-3M13 18H5v-6m-3 3 3-3 3 3`,server:`M3 2h18v8H3zM3 14h18v8H3zM7 6h1m3 0h6M7 18h1m3 0h6M6 10v4m12-4v4`,shield:`m12 2 8 3v7c0 5-8 10-8 10S4 17 4 12V5zM8 11l3 3 5-6`,upload:`M12 16V3m-5 5 5-5 5 5M4 15v6h16v-6`};function Ve(e){return ze.includes(e)?e:`parcel`}var He={"courier-button":ye,"courier-icon":class extends q{constructor(...e){super(...e),this.name=`parcel`,this.label=``}static{this.properties={name:{type:String},label:{type:String}}}static{this.styles=o`
    :host {
      display: inline-flex;
      width: 1.5rem;
      height: 1.5rem;
      color: currentColor;
    }
    svg { width: 100%; height: 100%; }
  `}render(){let e=Ve(this.name);return I`<svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="1.7"
      stroke-linecap="square"
      stroke-linejoin="miter"
      role=${this.label?`img`:`presentation`}
      aria-hidden=${this.label?`false`:`true`}
      aria-label=${this.label||void 0}
    ><path d=${Be[e]}></path></svg>`}},"courier-locale-selector":Ee,"courier-panel":De,"courier-progress":ke,"courier-theme-selector":Re};function Ue(e=customElements){for(let[t,n]of Object.entries(He))e.get(t)||e.define(t,n)}function Q(e,t=globalThis.location.pathname){return`${t.endsWith(`/`)?t:`${t}/`}api/v1/${e}`}async function $(e){if(!e.ok){let t=Error(`Courier administration request failed (${e.status})`);throw t.name=e.status===409?`ConflictError`:`RequestError`,t}if(e.status!==204)return e.json()}async function We(e=globalThis.fetch){return $(await e(Q(`servers`),{credentials:`same-origin`}))}async function Ge(e,t,n=globalThis.fetch){await $(await n(Q(`deliveries/${e.id}/policy`),{method:`PUT`,credentials:`same-origin`,headers:{"Content-Type":`application/json`},body:JSON.stringify({expectedVersion:e.policy.version,policy:t})}))}async function Ke(e,t,n=globalThis.fetch){await $(await n(Q(`${e}/${t}/stop`),{method:`POST`,credentials:`same-origin`}))}function qe(e,t=e=>new EventSource(e)){let n=t(Q(`events`));return n.addEventListener(`snapshot`,t=>e(JSON.parse(t.data))),()=>n.close()}var Je={en:{title:`Courier administration`,refresh:`Refresh`,retry:`Retry`,loading:`Loading servers…`,empty:`No Courier data servers found.`,failed:`Administration data is temporarily unavailable.`,conflict:`The policy changed elsewhere. Refresh before editing again.`,live:`Live`,unreachable:`Unreachable`,stopServer:`Stop server`,stopDelivery:`Stop delivery`,source:`Source`,destination:`Destination`,transferred:`Confirmed bytes`,authentication:`Authentication`,attempts:`Authentication attempts`,failAction:`Failure action`,noUi:`Hide delivery UI`,save:`Save policy`},ru:{title:`Управление Courier`,refresh:`Обновить`,retry:`Повторить`,loading:`Загрузка серверов…`,empty:`Серверы данных Courier не найдены.`,failed:`Данные управления временно недоступны.`,conflict:`Политика была изменена. Обновите данные перед повторным редактированием.`,live:`Работает`,unreachable:`Недоступен`,stopServer:`Остановить сервер`,stopDelivery:`Остановить доставку`,source:`Источник`,destination:`Назначение`,transferred:`Подтверждено байт`,authentication:`Аутентификация`,attempts:`Попытки аутентификации`,failAction:`Действие при ошибке`,noUi:`Скрыть интерфейс доставки`,save:`Сохранить политику`}};function Ye(e,t){return Je[e][t]}Ue();var Xe=class extends q{constructor(...e){super(...e),this.locale=Te(),this.failed=!1,this.conflict=!1}static{this.properties={locale:{state:!0},snapshot:{state:!0},failed:{state:!0},conflict:{state:!0}}}static{this.styles=o`
    :host { box-sizing: border-box; display: block; min-height: 100vh; padding: clamp(1rem, 4vw, 3rem); background: var(--courier-color-canvas); color: var(--courier-color-text); font-family: var(--courier-font-sans); }
    main, courier-panel { min-width: 0; }
    main { width: min(70rem, 100%); margin: 0 auto; display: grid; gap: 1rem; }
    header, nav, .row, .actions { display: flex; gap: .75rem; align-items: center; justify-content: space-between; flex-wrap: wrap; }
    section { display: grid; gap: .75rem; }
    article { padding: 1rem; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-medium); display: grid; gap: .75rem; }
    dl { display: grid; grid-template-columns: max-content 1fr; gap: .35rem .75rem; margin: 0; }
    dt { color: var(--courier-color-muted); }
    dd { margin: 0; overflow-wrap: anywhere; }
    form { display: grid; grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr)); gap: .75rem; align-items: end; }
    label { display: grid; gap: .25rem; }
    input, select { box-sizing: border-box; min-width: 0; max-width: 100%; min-height: 2.75rem; padding: 0 .5rem; }
    .status { text-transform: uppercase; letter-spacing: .08em; font-size: .75rem; }
    .status.live { color: var(--courier-color-accent); }
    h1, h2, p, strong { overflow-wrap: anywhere; }
    @media (max-width: 38rem) { dl { grid-template-columns: 1fr; } dt { margin-top: .35rem; } }
  `}connectedCallback(){super.connectedCallback(),this.theme=Le(),this.unsubscribe=qe(e=>{this.snapshot=e,this.failed=!1}),this.refresh()}disconnectedCallback(){this.unsubscribe?.(),this.theme?.destroy(),super.disconnectedCallback()}async refresh(){this.failed=!1,this.conflict=!1;try{this.snapshot=await We()}catch{this.failed=!0}}setLocale(e){this.locale=e.detail}async stop(e,t){try{await Ke(e,t),await this.refresh()}catch{this.failed=!0}}async save(e,t){e.preventDefault();let n=new FormData(e.currentTarget),r={...t.policy,version:t.policy.version+1,auth:String(n.get(`auth`)),authAttempts:Number(n.get(`attempts`)),authFailAction:String(n.get(`failAction`)),noUi:n.get(`noUi`)===`on`};this.failed=!1,this.conflict=!1;try{await Ge(t,r),await this.refresh()}catch(e){this.conflict=e instanceof Error&&e.name===`ConflictError`,this.failed=!this.conflict}}t(e){return Ye(this.locale,e)}delivery(e){return I`
      <article>
        <div class="row"><strong>${e.id}</strong><courier-button @click=${()=>this.stop(`deliveries`,e.id)}>${this.t(`stopDelivery`)}</courier-button></div>
        <dl>
          <dt>${this.t(`source`)}</dt><dd>${e.source||`<unavailable>`}</dd>
          <dt>${this.t(`destination`)}</dt><dd>${e.destination||`<unavailable>`}</dd>
          <dt>${this.t(`transferred`)}</dt><dd>${e.counters.confirmed}</dd>
        </dl>
        <form @submit=${t=>this.save(t,e)}>
          <label>${this.t(`authentication`)}<select name="auth"><option selected=${e.policy.auth===`none`}>none</option><option selected=${e.policy.auth===`basic`}>basic</option><option selected=${e.policy.auth===`password`}>password</option></select></label>
          <label>${this.t(`attempts`)}<input name="attempts" type="number" min="1" .value=${String(e.policy.authAttempts)}></label>
          <label>${this.t(`failAction`)}<select name="failAction"><option selected=${e.policy.authFailAction===`ban`}>ban</option><option selected=${e.policy.authFailAction===`stop`}>stop</option></select></label>
          <label><input name="noUi" type="checkbox" ?checked=${e.policy.noUi}> ${this.t(`noUi`)}</label>
          <courier-button type="submit">${this.t(`save`)}</courier-button>
        </form>
      </article>
    `}server(e){return I`
      <courier-panel>
        <section>
          <div class="row"><h2>${e.id}</h2><span class="status ${e.status}">${this.t(e.status)}</span></div>
          <div class="actions"><span>${e.bind}</span><courier-button @click=${()=>this.stop(`servers`,e.id)}>${this.t(`stopServer`)}</courier-button></div>
          ${e.deliveries.map(e=>this.delivery(e))}
        </section>
      </courier-panel>
    `}render(){let e=this.snapshot?.servers??[];return I`
      <main>
        <header>
          <h1>${this.t(`title`)}</h1>
          <nav><courier-button @click=${this.refresh}>${this.t(`refresh`)}</courier-button><courier-theme-selector></courier-theme-selector><courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector></nav>
        </header>
        ${this.failed?I`<courier-panel><p role="alert">${this.t(`failed`)}</p><courier-button @click=${this.refresh}>${this.t(`retry`)}</courier-button></courier-panel>`:R}
        ${this.conflict?I`<courier-panel><p role="alert">${this.t(`conflict`)}</p><courier-button @click=${this.refresh}>${this.t(`refresh`)}</courier-button></courier-panel>`:R}
        ${this.snapshot?e.length===0?I`<courier-panel>${this.t(`empty`)}</courier-panel>`:e.map(e=>this.server(e)):this.failed?R:I`<courier-panel>${this.t(`loading`)}</courier-panel>`}
      </main>
    `}};customElements.get(`courier-admin-app`)||customElements.define(`courier-admin-app`,Xe);