(function(){let e=document.createElement(`link`).relList;if(e&&e.supports&&e.supports(`modulepreload`))return;for(let e of document.querySelectorAll(`link[rel="modulepreload"]`))n(e);new MutationObserver(e=>{for(let t of e)if(t.type===`childList`)for(let e of t.addedNodes)e.tagName===`LINK`&&e.rel===`modulepreload`&&n(e)}).observe(document,{childList:!0,subtree:!0});function t(e){let t={};return e.integrity&&(t.integrity=e.integrity),e.referrerPolicy&&(t.referrerPolicy=e.referrerPolicy),t.credentials=e.crossOrigin===`use-credentials`?`include`:e.crossOrigin===`anonymous`?`omit`:`same-origin`,t}function n(e){if(e.ep)return;e.ep=!0;let n=t(e);fetch(e.href,n)}})();var e=globalThis,t=e.ShadowRoot&&(e.ShadyCSS===void 0||e.ShadyCSS.nativeShadow)&&`adoptedStyleSheets`in Document.prototype&&`replace`in CSSStyleSheet.prototype,n=Symbol(),r=new WeakMap,i=class{constructor(e,t,r){if(this._$cssResult$=!0,r!==n)throw Error("CSSResult is not constructable. Use `unsafeCSS` or `css` instead.");this.cssText=e,this.t=t}get styleSheet(){let e=this.o,n=this.t;if(t&&e===void 0){let t=n!==void 0&&n.length===1;t&&(e=r.get(n)),e===void 0&&((this.o=e=new CSSStyleSheet).replaceSync(this.cssText),t&&r.set(n,e))}return e}toString(){return this.cssText}},a=e=>new i(typeof e==`string`?e:e+``,void 0,n),o=(e,...t)=>new i(e.length===1?e[0]:t.reduce((t,n,r)=>t+(e=>{if(!0===e._$cssResult$)return e.cssText;if(typeof e==`number`)return e;throw Error(`Value passed to 'css' function must be a 'css' function result: `+e+`. Use 'unsafeCSS' to pass non-literal values, but take care to ensure page security.`)})(n)+e[r+1],e[0]),e,n),s=(n,r)=>{if(t)n.adoptedStyleSheets=r.map(e=>e instanceof CSSStyleSheet?e:e.styleSheet);else for(let t of r){let r=document.createElement(`style`),i=e.litNonce;i!==void 0&&r.setAttribute(`nonce`,i),r.textContent=t.cssText,n.appendChild(r)}},c=t?e=>e:e=>e instanceof CSSStyleSheet?(e=>{let t=``;for(let n of e.cssRules)t+=n.cssText;return a(t)})(e):e,{is:l,defineProperty:u,getOwnPropertyDescriptor:d,getOwnPropertyNames:ee,getOwnPropertySymbols:te,getPrototypeOf:ne}=Object,f=globalThis,p=f.trustedTypes,re=p?p.emptyScript:``,ie=f.reactiveElementPolyfillSupport,m=(e,t)=>e,h={toAttribute(e,t){switch(t){case Boolean:e=e?re:null;break;case Object:case Array:e=e==null?e:JSON.stringify(e)}return e},fromAttribute(e,t){let n=e;switch(t){case Boolean:n=e!==null;break;case Number:n=e===null?null:Number(e);break;case Object:case Array:try{n=JSON.parse(e)}catch{n=null}}return n}},g=(e,t)=>!l(e,t),_={attribute:!0,type:String,converter:h,reflect:!1,useDefault:!1,hasChanged:g};Symbol.metadata??=Symbol(`metadata`),f.litPropertyMetadata??=new WeakMap;var v=class extends HTMLElement{static addInitializer(e){this._$Ei(),(this.l??=[]).push(e)}static get observedAttributes(){return this.finalize(),this._$Eh&&[...this._$Eh.keys()]}static createProperty(e,t=_){if(t.state&&(t.attribute=!1),this._$Ei(),this.prototype.hasOwnProperty(e)&&((t=Object.create(t)).wrapped=!0),this.elementProperties.set(e,t),!t.noAccessor){let n=Symbol(),r=this.getPropertyDescriptor(e,n,t);r!==void 0&&u(this.prototype,e,r)}}static getPropertyDescriptor(e,t,n){let{get:r,set:i}=d(this.prototype,e)??{get(){return this[t]},set(e){this[t]=e}};return{get:r,set(t){let a=r?.call(this);i?.call(this,t),this.requestUpdate(e,a,n)},configurable:!0,enumerable:!0}}static getPropertyOptions(e){return this.elementProperties.get(e)??_}static _$Ei(){if(this.hasOwnProperty(m(`elementProperties`)))return;let e=ne(this);e.finalize(),e.l!==void 0&&(this.l=[...e.l]),this.elementProperties=new Map(e.elementProperties)}static finalize(){if(this.hasOwnProperty(m(`finalized`)))return;if(this.finalized=!0,this._$Ei(),this.hasOwnProperty(m(`properties`))){let e=this.properties,t=[...ee(e),...te(e)];for(let n of t)this.createProperty(n,e[n])}let e=this[Symbol.metadata];if(e!==null){let t=litPropertyMetadata.get(e);if(t!==void 0)for(let[e,n]of t)this.elementProperties.set(e,n)}this._$Eh=new Map;for(let[e,t]of this.elementProperties){let n=this._$Eu(e,t);n!==void 0&&this._$Eh.set(n,e)}this.elementStyles=this.finalizeStyles(this.styles)}static finalizeStyles(e){let t=[];if(Array.isArray(e)){let n=new Set(e.flat(1/0).reverse());for(let e of n)t.unshift(c(e))}else e!==void 0&&t.push(c(e));return t}static _$Eu(e,t){let n=t.attribute;return!1===n?void 0:typeof n==`string`?n:typeof e==`string`?e.toLowerCase():void 0}constructor(){super(),this._$Ep=void 0,this.isUpdatePending=!1,this.hasUpdated=!1,this._$Em=null,this._$Ev()}_$Ev(){this._$ES=new Promise(e=>this.enableUpdating=e),this._$AL=new Map,this._$E_(),this.requestUpdate(),this.constructor.l?.forEach(e=>e(this))}addController(e){(this._$EO??=new Set).add(e),this.renderRoot!==void 0&&this.isConnected&&e.hostConnected?.()}removeController(e){this._$EO?.delete(e)}_$E_(){let e=new Map,t=this.constructor.elementProperties;for(let n of t.keys())this.hasOwnProperty(n)&&(e.set(n,this[n]),delete this[n]);e.size>0&&(this._$Ep=e)}createRenderRoot(){let e=this.shadowRoot??this.attachShadow(this.constructor.shadowRootOptions);return s(e,this.constructor.elementStyles),e}connectedCallback(){this.renderRoot??=this.createRenderRoot(),this.enableUpdating(!0),this._$EO?.forEach(e=>e.hostConnected?.())}enableUpdating(e){}disconnectedCallback(){this._$EO?.forEach(e=>e.hostDisconnected?.())}attributeChangedCallback(e,t,n){this._$AK(e,n)}_$ET(e,t){let n=this.constructor.elementProperties.get(e),r=this.constructor._$Eu(e,n);if(r!==void 0&&!0===n.reflect){let i=(n.converter?.toAttribute===void 0?h:n.converter).toAttribute(t,n.type);this._$Em=e,i==null?this.removeAttribute(r):this.setAttribute(r,i),this._$Em=null}}_$AK(e,t){let n=this.constructor,r=n._$Eh.get(e);if(r!==void 0&&this._$Em!==r){let e=n.getPropertyOptions(r),i=typeof e.converter==`function`?{fromAttribute:e.converter}:e.converter?.fromAttribute===void 0?h:e.converter;this._$Em=r;let a=i.fromAttribute(t,e.type);this[r]=a??this._$Ej?.get(r)??a,this._$Em=null}}requestUpdate(e,t,n,r=!1,i){if(e!==void 0){let a=this.constructor;if(!1===r&&(i=this[e]),n??=a.getPropertyOptions(e),!((n.hasChanged??g)(i,t)||n.useDefault&&n.reflect&&i===this._$Ej?.get(e)&&!this.hasAttribute(a._$Eu(e,n))))return;this.C(e,t,n)}!1===this.isUpdatePending&&(this._$ES=this._$EP())}C(e,t,{useDefault:n,reflect:r,wrapped:i},a){n&&!(this._$Ej??=new Map).has(e)&&(this._$Ej.set(e,a??t??this[e]),!0!==i||a!==void 0)||(this._$AL.has(e)||(this.hasUpdated||n||(t=void 0),this._$AL.set(e,t)),!0===r&&this._$Em!==e&&(this._$Eq??=new Set).add(e))}async _$EP(){this.isUpdatePending=!0;try{await this._$ES}catch(e){Promise.reject(e)}let e=this.scheduleUpdate();return e!=null&&await e,!this.isUpdatePending}scheduleUpdate(){return this.performUpdate()}performUpdate(){if(!this.isUpdatePending)return;if(!this.hasUpdated){if(this.renderRoot??=this.createRenderRoot(),this._$Ep){for(let[e,t]of this._$Ep)this[e]=t;this._$Ep=void 0}let e=this.constructor.elementProperties;if(e.size>0)for(let[t,n]of e){let{wrapped:e}=n,r=this[t];!0!==e||this._$AL.has(t)||r===void 0||this.C(t,void 0,n,r)}}let e=!1,t=this._$AL;try{e=this.shouldUpdate(t),e?(this.willUpdate(t),this._$EO?.forEach(e=>e.hostUpdate?.()),this.update(t)):this._$EM()}catch(t){throw e=!1,this._$EM(),t}e&&this._$AE(t)}willUpdate(e){}_$AE(e){this._$EO?.forEach(e=>e.hostUpdated?.()),this.hasUpdated||(this.hasUpdated=!0,this.firstUpdated(e)),this.updated(e)}_$EM(){this._$AL=new Map,this.isUpdatePending=!1}get updateComplete(){return this.getUpdateComplete()}getUpdateComplete(){return this._$ES}shouldUpdate(e){return!0}update(e){this._$Eq&&=this._$Eq.forEach(e=>this._$ET(e,this[e])),this._$EM()}updated(e){}firstUpdated(e){}};v.elementStyles=[],v.shadowRootOptions={mode:`open`},v[m(`elementProperties`)]=new Map,v[m(`finalized`)]=new Map,ie?.({ReactiveElement:v}),(f.reactiveElementVersions??=[]).push(`2.1.2`);var y=globalThis,b=e=>e,x=y.trustedTypes,S=x?x.createPolicy(`lit-html`,{createHTML:e=>e}):void 0,C=`$lit$`,w=`lit$${Math.random().toFixed(9).slice(2)}$`,ae=`?`+w,oe=`<${ae}>`,T=document,E=()=>T.createComment(``),D=e=>e===null||typeof e!=`object`&&typeof e!=`function`,O=Array.isArray,se=e=>O(e)||typeof e?.[Symbol.iterator]==`function`,k=`[ 	
\f\r]`,A=/<(?:(!--|\/[^a-zA-Z])|(\/?[a-zA-Z][^>\s]*)|(\/?$))/g,j=/-->/g,M=/>/g,N=RegExp(`>|${k}(?:([^\\s"'>=/]+)(${k}*=${k}*(?:[^ \t\n\f\r"'\`<>=]|("|')|))|$)`,`g`),P=/'/g,F=/"/g,I=/^(?:script|style|textarea|title)$/i,L=(e=>(t,...n)=>({_$litType$:e,strings:t,values:n}))(1),R=Symbol.for(`lit-noChange`),z=Symbol.for(`lit-nothing`),B=new WeakMap,V=T.createTreeWalker(T,129);function H(e,t){if(!O(e)||!e.hasOwnProperty(`raw`))throw Error(`invalid template strings array`);return S===void 0?t:S.createHTML(t)}var ce=(e,t)=>{let n=e.length-1,r=[],i,a=t===2?`<svg>`:t===3?`<math>`:``,o=A;for(let t=0;t<n;t++){let n=e[t],s,c,l=-1,u=0;for(;u<n.length&&(o.lastIndex=u,c=o.exec(n),c!==null);)u=o.lastIndex,o===A?c[1]===`!--`?o=j:c[1]===void 0?c[2]===void 0?c[3]!==void 0&&(o=N):(I.test(c[2])&&(i=RegExp(`</`+c[2],`g`)),o=N):o=M:o===N?c[0]===`>`?(o=i??A,l=-1):c[1]===void 0?l=-2:(l=o.lastIndex-c[2].length,s=c[1],o=c[3]===void 0?N:c[3]===`"`?F:P):o===F||o===P?o=N:o===j||o===M?o=A:(o=N,i=void 0);let d=o===N&&e[t+1].startsWith(`/>`)?` `:``;a+=o===A?n+oe:l>=0?(r.push(s),n.slice(0,l)+C+n.slice(l)+w+d):n+w+(l===-2?t:d)}return[H(e,a+(e[n]||`<?>`)+(t===2?`</svg>`:t===3?`</math>`:``)),r]},U=class e{constructor({strings:t,_$litType$:n},r){let i;this.parts=[];let a=0,o=0,s=t.length-1,c=this.parts,[l,u]=ce(t,n);if(this.el=e.createElement(l,r),V.currentNode=this.el.content,n===2||n===3){let e=this.el.content.firstChild;e.replaceWith(...e.childNodes)}for(;(i=V.nextNode())!==null&&c.length<s;){if(i.nodeType===1){if(i.hasAttributes())for(let e of i.getAttributeNames())if(e.endsWith(C)){let t=u[o++],n=i.getAttribute(e).split(w),r=/([.?@])?(.*)/.exec(t);c.push({type:1,index:a,name:r[2],strings:n,ctor:r[1]===`.`?ue:r[1]===`?`?de:r[1]===`@`?fe:K}),i.removeAttribute(e)}else e.startsWith(w)&&(c.push({type:6,index:a}),i.removeAttribute(e));if(I.test(i.tagName)){let e=i.textContent.split(w),t=e.length-1;if(t>0){i.textContent=x?x.emptyScript:``;for(let n=0;n<t;n++)i.append(e[n],E()),V.nextNode(),c.push({type:2,index:++a});i.append(e[t],E())}}}else if(i.nodeType===8){if(i.data===ae)c.push({type:2,index:a});else{let e=-1;for(;(e=i.data.indexOf(w,e+1))!==-1;)c.push({type:7,index:a}),e+=w.length-1}}a++}}static createElement(e,t){let n=T.createElement(`template`);return n.innerHTML=e,n}};function W(e,t,n=e,r){if(t===R)return t;let i=r===void 0?n._$Cl:n._$Co?.[r],a=D(t)?void 0:t._$litDirective$;return i?.constructor!==a&&(i?._$AO?.(!1),a===void 0?i=void 0:(i=new a(e),i._$AT(e,n,r)),r===void 0?n._$Cl=i:(n._$Co??=[])[r]=i),i!==void 0&&(t=W(e,i._$AS(e,t.values),i,r)),t}var le=class{constructor(e,t){this._$AV=[],this._$AN=void 0,this._$AD=e,this._$AM=t}get parentNode(){return this._$AM.parentNode}get _$AU(){return this._$AM._$AU}u(e){let{el:{content:t},parts:n}=this._$AD,r=(e?.creationScope??T).importNode(t,!0);V.currentNode=r;let i=V.nextNode(),a=0,o=0,s=n[0];for(;s!==void 0;){if(a===s.index){let t;s.type===2?t=new G(i,i.nextSibling,this,e):s.type===1?t=new s.ctor(i,s.name,s.strings,this,e):s.type===6&&(t=new pe(i,this,e)),this._$AV.push(t),s=n[++o]}a!==s?.index&&(i=V.nextNode(),a++)}return V.currentNode=T,r}p(e){let t=0;for(let n of this._$AV)n!==void 0&&(n.strings===void 0?n._$AI(e[t]):(n._$AI(e,n,t),t+=n.strings.length-2)),t++}},G=class e{get _$AU(){return this._$AM?._$AU??this._$Cv}constructor(e,t,n,r){this.type=2,this._$AH=z,this._$AN=void 0,this._$AA=e,this._$AB=t,this._$AM=n,this.options=r,this._$Cv=r?.isConnected??!0}get parentNode(){let e=this._$AA.parentNode,t=this._$AM;return t!==void 0&&e?.nodeType===11&&(e=t.parentNode),e}get startNode(){return this._$AA}get endNode(){return this._$AB}_$AI(e,t=this){e=W(this,e,t),D(e)?e===z||e==null||e===``?(this._$AH!==z&&this._$AR(),this._$AH=z):e!==this._$AH&&e!==R&&this._(e):e._$litType$===void 0?e.nodeType===void 0?se(e)?this.k(e):this._(e):this.T(e):this.$(e)}O(e){return this._$AA.parentNode.insertBefore(e,this._$AB)}T(e){this._$AH!==e&&(this._$AR(),this._$AH=this.O(e))}_(e){this._$AH!==z&&D(this._$AH)?this._$AA.nextSibling.data=e:this.T(T.createTextNode(e)),this._$AH=e}$(e){let{values:t,_$litType$:n}=e,r=typeof n==`number`?this._$AC(e):(n.el===void 0&&(n.el=U.createElement(H(n.h,n.h[0]),this.options)),n);if(this._$AH?._$AD===r)this._$AH.p(t);else{let e=new le(r,this),n=e.u(this.options);e.p(t),this.T(n),this._$AH=e}}_$AC(e){let t=B.get(e.strings);return t===void 0&&B.set(e.strings,t=new U(e)),t}k(t){O(this._$AH)||(this._$AH=[],this._$AR());let n=this._$AH,r,i=0;for(let a of t)i===n.length?n.push(r=new e(this.O(E()),this.O(E()),this,this.options)):r=n[i],r._$AI(a),i++;i<n.length&&(this._$AR(r&&r._$AB.nextSibling,i),n.length=i)}_$AR(e=this._$AA.nextSibling,t){for(this._$AP?.(!1,!0,t);e!==this._$AB;){let t=b(e).nextSibling;b(e).remove(),e=t}}setConnected(e){this._$AM===void 0&&(this._$Cv=e,this._$AP?.(e))}},K=class{get tagName(){return this.element.tagName}get _$AU(){return this._$AM._$AU}constructor(e,t,n,r,i){this.type=1,this._$AH=z,this._$AN=void 0,this.element=e,this.name=t,this._$AM=r,this.options=i,n.length>2||n[0]!==``||n[1]!==``?(this._$AH=Array(n.length-1).fill(new String),this.strings=n):this._$AH=z}_$AI(e,t=this,n,r){let i=this.strings,a=!1;if(i===void 0)e=W(this,e,t,0),a=!D(e)||e!==this._$AH&&e!==R,a&&(this._$AH=e);else{let r=e,o,s;for(e=i[0],o=0;o<i.length-1;o++)s=W(this,r[n+o],t,o),s===R&&(s=this._$AH[o]),a||=!D(s)||s!==this._$AH[o],s===z?e=z:e!==z&&(e+=(s??``)+i[o+1]),this._$AH[o]=s}a&&!r&&this.j(e)}j(e){e===z?this.element.removeAttribute(this.name):this.element.setAttribute(this.name,e??``)}},ue=class extends K{constructor(){super(...arguments),this.type=3}j(e){this.element[this.name]=e===z?void 0:e}},de=class extends K{constructor(){super(...arguments),this.type=4}j(e){this.element.toggleAttribute(this.name,!!e&&e!==z)}},fe=class extends K{constructor(e,t,n,r,i){super(e,t,n,r,i),this.type=5}_$AI(e,t=this){if((e=W(this,e,t,0)??z)===R)return;let n=this._$AH,r=e===z&&n!==z||e.capture!==n.capture||e.once!==n.once||e.passive!==n.passive,i=e!==z&&(n===z||r);r&&this.element.removeEventListener(this.name,this,n),i&&this.element.addEventListener(this.name,this,e),this._$AH=e}handleEvent(e){typeof this._$AH==`function`?this._$AH.call(this.options?.host??this.element,e):this._$AH.handleEvent(e)}},pe=class{constructor(e,t,n){this.element=e,this.type=6,this._$AN=void 0,this._$AM=t,this.options=n}get _$AU(){return this._$AM._$AU}_$AI(e){W(this,e)}},me=y.litHtmlPolyfillSupport;me?.(U,G),(y.litHtmlVersions??=[]).push(`3.3.3`);var he=(e,t,n)=>{let r=n?.renderBefore??t,i=r._$litPart$;if(i===void 0){let e=n?.renderBefore??null;r._$litPart$=i=new G(t.insertBefore(E(),e),e,void 0,n??{})}return i._$AI(e),i},q=globalThis,J=class extends v{constructor(){super(...arguments),this.renderOptions={host:this},this._$Do=void 0}createRenderRoot(){let e=super.createRenderRoot();return this.renderOptions.renderBefore??=e.firstChild,e}update(e){let t=this.render();this.hasUpdated||(this.renderOptions.isConnected=this.isConnected),super.update(e),this._$Do=he(t,this.renderRoot,this.renderOptions)}connectedCallback(){super.connectedCallback(),this._$Do?.setConnected(!0)}disconnectedCallback(){super.disconnectedCallback(),this._$Do?.setConnected(!1)}render(){return R}};J._$litElement$=!0,J.finalized=!0,q.litElementHydrateSupport?.({LitElement:J});var ge=q.litElementPolyfillSupport;ge?.({LitElement:J}),(q.litElementVersions??=[]).push(`4.2.2`);var _e=o`
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
    outline: 3px solid var(--courier-beak, #ff8758);
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
`;var ve=o`
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
    outline: 3px solid var(--courier-beak, #ff8758);
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
    outline: 3px solid var(--courier-beak, #ff8758);
    outline-offset: 2px;
  }
`,ye=class extends J{constructor(...e){super(...e),this.disabled=!1,this.type=`button`,this.variant=`secondary`}static{this.properties={disabled:{type:Boolean,reflect:!0},type:{type:String,reflect:!0},variant:{type:String,reflect:!0}}}static{this.styles=[_e,o`
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
  `]}render(){return L`<button type=${this.type} ?disabled=${this.disabled}><slot></slot></button>`}},be=class extends J{constructor(...e){super(...e),this.product=``}static{this.properties={product:{type:String}}}static{this.styles=o`
    :host { display: inline-flex; min-width: 0; color: var(--courier-color-text, #151714); font-family: var(--courier-font-sans, sans-serif); }
    .lockup { display: inline-flex; min-width: 0; align-items: center; gap: 0.625rem; color: inherit; }
    svg { flex: 0 0 auto; width: 2.25rem; height: 2.25rem; }
    .words { display: grid; line-height: 1; }
    strong { font-family: var(--courier-font-display, sans-serif); font-size: 1.125rem; font-weight: 850; letter-spacing: -0.04em; }
    small { margin-top: 0.25rem; color: var(--courier-color-muted, #596054); font-family: var(--courier-font-mono, monospace); font-size: 0.625rem; font-weight: 700; letter-spacing: 0.09em; text-transform: uppercase; }
  `}render(){return L`<span class="lockup">
      <svg viewBox="0 0 40 40" aria-hidden="true">
        <rect x="1" y="1" width="38" height="38" rx="4" fill="var(--courier-graphite-900, #151714)" stroke="var(--courier-color-border-strong, #8e9587)"></rect>
        <path d="M10 30c-1.5-3.2-1.4-7 .3-10.2C12.5 15.5 16.4 12 21.5 12c3.8 0 7 1.8 8.8 4.7l-2.5 3.6c1 2.8.8 6.4-.7 9.7H10z" fill="var(--courier-paper-50, #f3f4e9)"></path>
        <path d="m29 15 8 3.1-8 4z" fill="var(--courier-beak, #ff8758)"></path>
        <circle cx="24" cy="16" r="2.2" fill="var(--courier-graphite-900, #151714)"></circle>
        <circle cx="24" cy="16" r=".8" fill="var(--courier-color-accent, #d4ff45)"></circle>
        <path d="M10 25h18.2c.1 1.7-.3 3.5-1.1 5H10c-.6-1.6-.8-3.3-.7-5z" fill="var(--courier-color-accent, #d4ff45)"></path>
        <path d="M15 25v5m8-5v5" stroke="var(--courier-graphite-900, #151714)" stroke-width="1.2"></path>
      </svg>
      <span class="words"><strong>Courier</strong><small>${this.product}</small></span>
    </span>`}},xe=class extends J{constructor(...e){super(...e),this.alt=``,this.eager=!1,this.source=``}static{this.properties={alt:{type:String},eager:{type:Boolean},source:{type:String}}}static{this.styles=o`
    :host { display: block; }
    img { display: block; width: 100%; height: auto; filter: drop-shadow(0 1.5rem 2rem rgb(16 18 15 / 0.18)); }
  `}render(){return L`<img part="image" src=${this.source} alt=${this.alt} decoding="async" loading=${this.eager?`eager`:`lazy`} fetchpriority=${this.eager?`high`:`auto`}>`}},Se=class extends J{constructor(...e){super(...e),this.tone=`neutral`}static{this.properties={tone:{type:String,reflect:!0}}}static{this.styles=o`
    :host {
      display: inline-flex;
      width: max-content;
      align-items: center;
      gap: 0.45rem;
      color: var(--courier-color-muted, #596054);
      font-family: var(--courier-font-mono, monospace);
      font-size: 0.6875rem;
      font-weight: 750;
      letter-spacing: 0.075em;
      text-transform: uppercase;
    }
    i { width: 0.5rem; height: 0.5rem; border: 1px solid currentColor; border-radius: 50%; background: currentColor; box-shadow: 0 0 0 3px color-mix(in srgb, currentColor 14%, transparent); }
    :host([tone="signal"]) { color: var(--courier-success, #76a51f); }
    :host([tone="warning"]) { color: var(--courier-warning, #c78300); }
    :host([tone="danger"]) { color: var(--courier-danger, #ff6b5f); }
  `}render(){return L`<i aria-hidden="true"></i><slot></slot>`}},Ce=class extends J{constructor(...e){super(...e),this.source=``,this.destination=``}static{this.properties={source:{type:String},destination:{type:String}}}static{this.styles=o`
    :host { display: grid; color: var(--courier-color-text, #151714); font-family: var(--courier-font-mono, monospace); }
    .route { display: grid; grid-template-columns: minmax(0, 1fr) minmax(3rem, 0.55fr) minmax(0, 1fr); align-items: center; gap: 0.65rem; }
    .node { overflow: hidden; padding: 0.65rem 0.75rem; border: 1px solid var(--courier-color-border, #c8cdbf); border-radius: var(--courier-radius-sm, 0.25rem); background: var(--courier-color-surface, #fafbf3); font-size: 0.75rem; text-overflow: ellipsis; white-space: nowrap; }
    .line { position: relative; height: 1px; color: var(--courier-color-border-strong, #8e9587); background: currentColor; }
    .line::before { content: ""; position: absolute; top: -0.2rem; left: 0; width: 0.45rem; height: 0.45rem; border-radius: 50%; background: var(--courier-color-accent, #d4ff45); }
    .line::after { content: ""; position: absolute; top: -0.22rem; right: 0; width: 0.4rem; height: 0.4rem; border-top: 1px solid currentColor; border-right: 1px solid currentColor; transform: rotate(45deg); }
  `}render(){return L`<div class="route"><span class="node">${this.source}</span><span class="line" aria-hidden="true"></span><span class="node">${this.destination}</span></div>`}},we={en:{"theme.label":`Theme`,"theme.system":`System`,"theme.light":`Light`,"theme.dark":`Dark`,"locale.label":`Language`,"locale.en":`English`,"locale.ru":`Russian`,"progress.label":`Delivery progress`,"action.cancel":`Cancel`,"action.close":`Close`},ru:{"theme.label":`Тема`,"theme.system":`Системная`,"theme.light":`Светлая`,"theme.dark":`Тёмная`,"locale.label":`Язык`,"locale.en":`Английский`,"locale.ru":`Русский`,"progress.label":`Ход доставки`,"action.cancel":`Отмена`,"action.close":`Закрыть`}},Te=Object.keys(we),Y=`courier.locale`;function X(e){if(!e)return;let t=e.toLowerCase().split(`-`)[0];return Te.includes(t)?t:void 0}function Ee(e){for(let t of e){let e=X(t);if(e)return e}return`en`}function De(e,t){if(e)try{let t=X(e.getItem(Y));if(t)return t}catch{}return Ee(t)}function Oe(e,t){if(e)try{e.setItem(Y,t)}catch{}}function Z(e,t){return we[X(e)??`en`][t]}function ke(){let e;try{e=globalThis.localStorage}catch{e=void 0}let t=globalThis.navigator?.languages??[];return De(e,t)}var Ae=class extends J{constructor(...e){super(...e),this.locale=`en`}static{this.properties={locale:{type:String}}}connectedCallback(){super.connectedCallback(),this.locale=ke()}change(e){let t=X(e.detail)??`en`;this.locale=t;let n;try{n=globalThis.localStorage}catch{n=void 0}Oe(n,t),this.dispatchEvent(new CustomEvent(`courier-locale-change`,{detail:t,bubbles:!0,composed:!0}))}render(){return L`<courier-segmented-control
      icon-only
      .label=${Z(this.locale,`locale.label`)}
      .value=${this.locale}
      .options=${[{value:`en`,label:Z(this.locale,`locale.en`),icon:`language-en`},{value:`ru`,label:Z(this.locale,`locale.ru`),icon:`language-ru`}]}
      data-storage-key=${Y}
      @courier-segment-change=${this.change}
    ></courier-segmented-control>`}},je=class extends J{constructor(...e){super(...e),this.heading=``}static{this.properties={heading:{type:String}}}static{this.styles=o`
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
  `}render(){return this.heading?L`<section aria-labelledby="courier-panel-heading">
      <h2 id="courier-panel-heading">${this.heading}</h2>
      <slot></slot>
    </section>`:L`<section><slot></slot></section>`}};function Me(e,t){return!Number.isFinite(e)||!Number.isFinite(t)||t<=0?0:Math.min(1,Math.max(0,e/t))}var Ne=class extends J{constructor(...e){super(...e),this.value=0,this.total=0,this.label=``,this.locale=`en`}static{this.properties={value:{type:Number},total:{type:Number},label:{type:String},locale:{type:String}}}static{this.styles=o`
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
  `}render(){let e=Me(this.value,this.total),t=this.label||Z(this.locale,`progress.label`),n=Number.isFinite(this.total)&&this.total>0?this.total:0;return L`<div
      class="track"
      role="progressbar"
      aria-label=${t}
      aria-valuemin="0"
      aria-valuemax=${n}
      aria-valuenow=${Number.isFinite(this.value)?Math.max(0,Math.min(this.value,n)):0}
    ><div class="fill" style=${`transform: scaleX(${e})`}></div></div>
    <output>${Math.round(e*100)}%</output>`}};function Pe(e,t,n){if(!(n<=0))switch(e){case`ArrowLeft`:case`ArrowUp`:return(t-1+n)%n;case`ArrowRight`:case`ArrowDown`:return(t+1)%n;case`Home`:return 0;case`End`:return n-1;default:return}}var Fe=class extends J{constructor(...e){super(...e),this.label=``,this.options=[],this.value=``,this.iconOnly=!1}static{this.properties={label:{type:String},options:{attribute:!1},value:{type:String},iconOnly:{type:Boolean,attribute:`icon-only`,reflect:!0}}}static{this.styles=o`
    :host { display: block; min-width: 0; color: var(--courier-color-text, #151714); font-family: var(--courier-font-sans, sans-serif); }
    fieldset { min-width: 0; margin: 0; padding: 0; border: 0; }
    legend { margin: 0 0 0.25rem; padding: 0; color: var(--courier-color-muted, #596054); font-family: var(--courier-font-mono, monospace); font-size: 0.625rem; font-weight: 750; letter-spacing: 0.08em; line-height: 1; text-transform: uppercase; }
    legend.sr-only { position: absolute; width: 1px; height: 1px; margin: -1px; padding: 0; overflow: hidden; clip: rect(0 0 0 0); clip-path: inset(50%); white-space: nowrap; }
    .segments { display: inline-grid; max-width: 100%; grid-auto-columns: minmax(0, auto); grid-auto-flow: column; gap: 2px; padding: 2px; border: 1px solid var(--courier-color-border, #c8cdbf); border-radius: var(--courier-radius-md, 0.625rem); background: color-mix(in srgb, var(--courier-color-field, #e7e9dc) 68%, transparent); box-shadow: inset 0 1px 2px rgb(16 18 15 / 0.07); }
    button { appearance: none; display: inline-flex; min-width: 0; min-height: 2.25rem; align-items: center; justify-content: center; gap: 0.35rem; padding: 0.45rem 0.68rem; overflow: hidden; border: 1px solid transparent; border-radius: calc(var(--courier-radius-md, 0.625rem) - 3px); color: var(--courier-color-muted, #596054); background: transparent; font: inherit; font-size: 0.75rem; font-weight: 780; line-height: 1; text-overflow: ellipsis; white-space: nowrap; cursor: pointer; transition: color var(--courier-duration, 160ms) var(--courier-ease, ease), background var(--courier-duration, 160ms) var(--courier-ease, ease), transform var(--courier-duration, 160ms) var(--courier-ease, ease); }
    :host([icon-only]) button { width: 2.25rem; padding: 0.45rem; }
    courier-icon { width: 1.05rem; height: 1.05rem; }
    button:hover { color: var(--courier-color-text, #151714); background: color-mix(in srgb, var(--courier-color-surface-raised, #fff) 72%, transparent); }
    button.selected { border-color: color-mix(in srgb, var(--courier-color-accent, #d4ff45) 64%, var(--courier-color-border, #c8cdbf)); color: var(--courier-color-accent-ink, #151714); background: var(--courier-color-accent, #d4ff45); box-shadow: 0 1px 0 rgb(16 18 15 / 0.12); }
    button:active { transform: translateY(1px); }
    button:focus-visible { position: relative; z-index: 1; outline: 3px solid var(--courier-beak, #ff8758); outline-offset: 2px; }
  `}activate(e){let t=e.currentTarget.dataset.value??``;t&&t!==this.value&&(this.value=t,this.dispatchEvent(new CustomEvent(`courier-segment-change`,{detail:t,bubbles:!0,composed:!0})))}move(e){let t=[...this.renderRoot.querySelectorAll(`button`)],n=t.indexOf(e.currentTarget),r=Pe(e.key,n,t.length);if(r===void 0)return;e.preventDefault();let i=t[r];i.focus(),i.click()}render(){return L`<fieldset>
      <legend class=${this.iconOnly?`sr-only`:``}>${this.label}</legend>
      <div class="segments" role="radiogroup" aria-label=${this.label}>
        ${this.options.map(e=>{let t=e.value===this.value;return L`<button
            type="button"
            role="radio"
            class=${t?`selected`:``}
            data-value=${e.value}
            aria-label=${e.label}
            title=${e.label}
            aria-checked=${String(t)}
            tabindex=${t?0:-1}
            @click=${this.activate}
            @keydown=${this.move}
          >${e.icon?L`<courier-icon name=${e.icon}></courier-icon>`:``}${this.iconOnly?``:L`<span>${e.label}</span>`}</button>`})}
      </div>
    </fieldset>`}},Ie=[`system`,`light`,`dark`],Le=`courier.theme`;function Re(e){return Ie.includes(e)?e:`system`}function ze(e,t){return e===`system`?t?.matches?`dark`:`light`:e}function Be(e){if(!e)return`system`;try{return Re(e.getItem(Le))}catch{return`system`}}function Ve(e,t){if(e)try{e.setItem(Le,t)}catch{}}var He=class{constructor(e,t,n,r){this.root=e,this.storage=t,this.media=n,this.onSystemChange=()=>this.apply(),this.preference=r??Be(t),this.media?.addEventListener(`change`,this.onSystemChange),this.apply()}set(e){this.preference=e,Ve(this.storage,e),this.apply()}destroy(){this.media?.removeEventListener(`change`,this.onSystemChange)}apply(){this.root.dataset.courierTheme=ze(this.preference,this.media),this.root.dataset.courierThemePreference=this.preference}};function Ue(){let e;try{e=globalThis.localStorage}catch{e=void 0}let t=globalThis.matchMedia?.(`(prefers-color-scheme: dark)`);return new He(document.documentElement,e,t)}var We=class extends J{constructor(...e){super(...e),this.preference=`system`,this.locale=`en`}static{this.properties={preference:{type:String},locale:{type:String}}}connectedCallback(){super.connectedCallback(),this.state=Ue(),this.preference=this.state.preference}disconnectedCallback(){this.state?.destroy(),super.disconnectedCallback()}change(e){let t=Re(e.detail);this.preference=t,this.state?.set(t),this.dispatchEvent(new CustomEvent(`courier-theme-change`,{detail:t,bubbles:!0,composed:!0}))}render(){return L`<courier-segmented-control
      icon-only
      .label=${Z(this.locale,`theme.label`)}
      .value=${this.preference}
      .options=${[{value:`system`,label:Z(this.locale,`theme.system`),icon:`system`},{value:`light`,label:Z(this.locale,`theme.light`),icon:`sun`},{value:`dark`,label:Z(this.locale,`theme.dark`),icon:`moon`}]}
      @courier-segment-change=${this.change}
    ></courier-segmented-control>`}},Ge=`apple.archive.copy.download.folder.github.homebrew.language-en.language-ru.linux.moon.npm.package.parcel.pnpm.receipt.retry.route.scoop.server.shield.sun.system.terminal.upload.windows.yarn`.split(`.`),Ke={apple:`M15 5c1-1 1-3 1-3-2 0-3 1-4 3m6 7c-1-2-2-3-4-3-1 0-2 1-3 1s-2-1-3-1c-3 0-5 3-5 6 0 4 3 8 5 8 1 0 2-1 3-1s2 1 3 1c2 0 4-3 5-6-2-1-3-2-3-4 0-2 1-3 2-4z`,archive:`M3 3h18v5H3zM5 8v13h14V8M9 12h6`,copy:`M8 3h13v13M3 8h13v13H3z`,download:`M12 3v13m-5-5 5 5 5-5M4 15v6h16v-6`,folder:`M3 6h7l2 3h9v12H3zM3 6V3h7l2 3h9v3`,github:`M9 19c-5 1-5-2-7-3m14 6v-3.6c0-1 .1-1.7-.4-2.2 3.2-.4 6.4-1.6 6.4-7.1 0-1.6-.6-3-1.7-4 .2-.5.7-2.3-.2-4.6 0 0-1.4-.5-4.7 1.7a16 16 0 0 0-8.6 0C6.4 1 5 1.5 5 1.5 4.1 3.8 4.6 5.6 4.8 6.1a7 7 0 0 0-1.7 4c0 5.5 3.2 6.7 6.4 7.1-.4.4-.8 1.1-.8 2.2V23`,homebrew:`M6 4h11l-1 15H8zM17 7h2a2 2 0 0 1 2 2v3a2 2 0 0 1-2 2h-2M5 22h13`,"language-en":`M3 4h18v14H9l-4 3v-3H3z`,"language-ru":`M3 4h18v14H9l-4 3v-3H3z`,linux:`M12 2c-3 0-4 3-4 6-2 2-3 5-3 8l3-1 1 5 3-2 3 2 1-5 3 1c0-3-1-6-3-8 0-3-1-6-4-6zM9 8h.01M15 8h.01M10 11h4`,moon:`M20 16a8 8 0 0 1-12-10 8 8 0 1 0 12 10z`,npm:`M2 6h20v12H2zM6 15V9h5v6m0-6h4v6m0-6h3v6`,package:`m3 7 9-5 9 5v11l-9 4-9-4zM3 7l9 5 9-5M12 12v10M8 4l9 5`,parcel:`m3 7 9-5 9 5v11l-9 4-9-4zM3 7l9 5 9-5M12 12v10M8 4l9 5v5`,pnpm:`M3 3h5v5H3zM10 3h5v5h-5zM17 3h4v5h-4zM3 10h5v5H3zm7 0h5v5h-5zm7 0h4v5h-4zM10 17h5v4h-5zm7 0h4v4h-4z`,receipt:`M5 2h14v20l-3-2-4 2-4-2-3 2zM8 7h8M8 11h8m-8 5 2 2 5-4`,retry:`M3 10a9 9 0 1 1 1 7M3 3v7h7M12 7v5l3 2`,route:`M2 3h6v6H2zM16 15h6v6h-6zM11 6h8v6m-3-3 3 3 3-3M13 18H5v-6m-3 3 3-3 3 3`,scoop:`M5 8h14l-2 13H7zM4 8h16M8 8V5a4 4 0 0 1 8 0v3`,server:`M3 2h18v8H3zM3 14h18v8H3zM7 6h1m3 0h6M7 18h1m3 0h6M6 10v4m12-4v4`,shield:`m12 2 8 3v7c0 5-8 10-8 10S4 17 4 12V5zM8 11l3 3 5-6`,sun:`M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8zM12 2v3m0 14v3M4.9 4.9 7 7m10 10 2.1 2.1M2 12h3m14 0h3M4.9 19.1 7 17M17 7l2.1-2.1`,system:`M3 4h18v13H3zM8 21h8M12 17v4`,terminal:`M4 6l5 6-5 6m7 0h9`,upload:`M12 16V3m-5 5 5-5 5 5M4 15v6h16v-6`,windows:`M3 4l8-1v8H3zm10-1 8-1v9h-8zM3 13h8v8l-8-1zm10 0h8v9l-8-1z`,yarn:`M12 3a9 9 0 1 0 9 9M8 16c4-1 7-4 9-8m-8 1c3 1 5 4 5 8m-5-5c-1-3 0-5 2-6`},qe={"language-en":`EN`,"language-ru":`RU`};function Je(e){return Ge.includes(e)?e:`parcel`}var Ye={"courier-brand":be,"courier-button":ye,"courier-icon":class extends J{constructor(...e){super(...e),this.name=`parcel`,this.label=``}static{this.properties={name:{type:String},label:{type:String}}}static{this.styles=o`
    :host {
      display: inline-flex;
      width: 1.5rem;
      height: 1.5rem;
      color: currentColor;
    }
    svg { width: 100%; height: 100%; }
    text { fill: currentColor; stroke: none; font-family: var(--courier-font-mono, monospace); font-size: 6px; font-weight: 850; letter-spacing: -0.04em; text-anchor: middle; }
  `}render(){let e=Je(this.name);return L`<svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="1.7"
      stroke-linecap="square"
      stroke-linejoin="miter"
      role=${this.label?`img`:`presentation`}
      aria-hidden=${this.label?`false`:`true`}
      aria-label=${this.label||void 0}
    ><path d=${Ke[e]}></path><text x="12" y="13.1">${qe[e]??``}</text></svg>`}},"courier-locale-selector":Ae,"courier-mascot":xe,"courier-panel":je,"courier-progress":Ne,"courier-route":Ce,"courier-segmented-control":Fe,"courier-status":Se,"courier-theme-selector":We};function Xe(e=customElements){for(let[t,n]of Object.entries(Ye))e.get(t)||e.define(t,n)}var Ze=`/assets/admin.webp`;function Q(e,t=globalThis.location.pathname){return`${t.endsWith(`/`)?t:`${t}/`}api/v1/${e}`}async function $(e){if(!e.ok){let t=Error(`Courier administration request failed (${e.status})`);throw t.name=e.status===409?`ConflictError`:`RequestError`,t}if(e.status!==204)return e.json()}async function Qe(e=globalThis.fetch){return $(await e(Q(`servers`),{credentials:`same-origin`}))}async function $e(e,t,n=globalThis.fetch){await $(await n(Q(`deliveries/${e.id}/policy`),{method:`PUT`,credentials:`same-origin`,headers:{"Content-Type":`application/json`},body:JSON.stringify({expectedVersion:e.policy.version,policy:t})}))}async function et(e,t,n=globalThis.fetch){await $(await n(Q(`${e}/${t}/stop`),{method:`POST`,credentials:`same-origin`}))}function tt(e,t=e=>new EventSource(e)){let n=t(Q(`events`));return n.addEventListener(`snapshot`,t=>e(JSON.parse(t.data))),()=>n.close()}var nt={en:{brandProduct:`Operations control`,eyebrow:`Local control plane`,title:`Delivery control`,intro:`Inspect live routes, confirmed volume, and delivery policy from one private local console.`,refresh:`Refresh registry`,retry:`Retry connection`,loading:`Checking the live registry…`,empty:`No live Courier servers are registered.`,failed:`Administration data is temporarily unavailable. Retry the local connection.`,conflict:`This delivery policy changed elsewhere. Refresh before editing again.`,live:`Live`,unreachable:`Unreachable`,stopServer:`Stop server`,stopDelivery:`Stop delivery`,source:`Source`,destination:`Destination`,transferred:`Confirmed bytes`,authentication:`Authentication`,attempts:`Authentication attempts`,failAction:`Failure action`,noUi:`Hide delivery UI`,save:`Apply policy`,serversMetric:`Live registry servers`,deliveriesMetric:`Active deliveries`,confirmedMetric:`Confirmed bytes`,server:`Server`,bind:`Bound address`,deliveries:`Delivery routes`,policy:`Delivery policy`,unavailable:`Unavailable`},ru:{brandProduct:`Операционный контроль`,eyebrow:`Локальный контур управления`,title:`Управление доставками`,intro:`Проверяйте активные маршруты, подтверждённый объём и правила доставки в одной приватной локальной консоли.`,refresh:`Обновить реестр`,retry:`Повторить подключение`,loading:`Проверка активного реестра…`,empty:`Активные серверы Courier не зарегистрированы.`,failed:`Данные управления временно недоступны. Повторите локальное подключение.`,conflict:`Правила этой доставки были изменены. Обновите данные перед повторным редактированием.`,live:`Работает`,unreachable:`Недоступен`,stopServer:`Остановить сервер`,stopDelivery:`Остановить доставку`,source:`Источник`,destination:`Назначение`,transferred:`Подтверждено байт`,authentication:`Аутентификация`,attempts:`Попытки аутентификации`,failAction:`Действие при ошибке`,noUi:`Скрыть интерфейс доставки`,save:`Применить правила`,serversMetric:`Серверы активного реестра`,deliveriesMetric:`Активные доставки`,confirmedMetric:`Подтверждено байт`,server:`Сервер`,bind:`Адрес привязки`,deliveries:`Маршруты доставки`,policy:`Правила доставки`,unavailable:`Недоступно`}};function rt(e,t){return nt[e][t]}Xe();var it=class extends J{constructor(...e){super(...e),this.locale=ke(),this.failed=!1,this.conflict=!1}static{this.properties={locale:{state:!0},snapshot:{state:!0},failed:{state:!0},conflict:{state:!0}}}static{this.styles=[ve,o`
    :host {
      display: block;
      min-height: 100vh;
      padding: 0 1rem 4rem;
      color: var(--courier-color-text);
      background-color: var(--courier-color-canvas);
      background-image: linear-gradient(var(--courier-color-grid) 1px, transparent 1px), linear-gradient(90deg, var(--courier-color-grid) 1px, transparent 1px);
      background-size: 2.5rem 2.5rem;
      font-family: var(--courier-font-sans);
    }
    main, courier-panel, article { min-width: 0; }
    main { width: min(78rem, 100%); margin: 0 auto; }
    header { display: flex; min-height: 5rem; align-items: center; justify-content: space-between; gap: 1rem; border-bottom: 1px solid var(--courier-color-border); }
    nav, .row, .actions { display: flex; gap: 0.75rem; align-items: center; justify-content: space-between; flex-wrap: wrap; }
    .workspace { display: grid; gap: 1.25rem; padding-top: clamp(2rem, 6vw, 5rem); }
    .page-head { display: grid; grid-template-columns: minmax(0, 1fr) auto; gap: 2rem; align-items: end; padding-bottom: 1.25rem; border-bottom: 1px solid var(--courier-color-border); }
    .page-head > div { display: grid; gap: 0.65rem; }
    .eyebrow, .label { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; font-weight: 750; letter-spacing: 0.08em; text-transform: uppercase; }
    h1, h2, h3, p, strong { margin: 0; overflow-wrap: anywhere; }
    h1 { font-family: var(--courier-font-display); font-size: clamp(2.5rem, 7vw, 5rem); font-weight: 830; letter-spacing: -0.065em; line-height: 0.92; }
    h2 { font-size: 1.2rem; letter-spacing: -0.03em; }
    h3 { font-size: 1rem; }
    p { line-height: 1.6; }
    .intro { max-width: 43rem; color: var(--courier-color-muted); }
    .metrics { display: grid; grid-template-columns: repeat(3, 1fr); border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-md); background: var(--courier-color-surface-raised); box-shadow: var(--courier-shadow); }
    .metric { display: grid; gap: 0.35rem; padding: 1.35rem; }
    .metric + .metric { border-left: 1px solid var(--courier-color-border); }
    .metric strong { font-family: var(--courier-font-display); font-size: 2rem; font-variant-numeric: tabular-nums; letter-spacing: -0.05em; }
    .metric span { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; letter-spacing: 0.075em; text-transform: uppercase; }
    .notice { display: grid; gap: 0.75rem; padding: 1rem; border: 1px solid var(--courier-color-border); border-left: 3px solid var(--courier-warning); border-radius: var(--courier-radius-sm); background: var(--courier-color-surface-raised); }
    .notice.error { border-left-color: var(--courier-danger); }
    .server-list { display: grid; gap: 1rem; }
    .server { display: grid; gap: 1rem; }
    .server-head { display: grid; grid-template-columns: minmax(0, 1fr) auto; gap: 1rem; padding-bottom: 1rem; border-bottom: 1px solid var(--courier-color-border); }
    .server-title { display: grid; gap: 0.4rem; }
    .server-id, .delivery-id, dd { font-family: var(--courier-font-mono); font-size: 0.8rem; font-variant-numeric: tabular-nums; }
    .bind { display: flex; align-items: center; gap: 0.5rem; color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.8rem; }
    .delivery-stack { display: grid; gap: 0.75rem; }
    article { display: grid; gap: 1rem; padding: clamp(1rem, 3vw, 1.5rem); border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-md); background: var(--courier-color-surface); }
    .delivery-head { align-items: start; }
    .delivery-head > div { display: grid; gap: 0.35rem; }
    .route { display: grid; gap: 0.75rem; }
    dl { display: grid; grid-template-columns: max-content minmax(0, 1fr); gap: 0.45rem 1rem; margin: 0; padding: 0.9rem 0; border-top: 1px solid var(--courier-color-border); border-bottom: 1px solid var(--courier-color-border); }
    dt { color: var(--courier-color-muted); font-size: 0.8rem; }
    dd { margin: 0; overflow-wrap: anywhere; }
    form { display: grid; grid-template-columns: repeat(4, minmax(9rem, 1fr)) auto; gap: 0.75rem; align-items: end; }
    label { display: grid; gap: 0.35rem; color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; font-weight: 700; letter-spacing: 0.04em; }
    label.checkbox { grid-template-columns: auto 1fr; align-items: center; align-content: center; min-height: 2.75rem; }
    .state-brief { display: grid; grid-template-columns: minmax(0, 1fr) minmax(18rem, 0.7fr); min-height: 20rem; overflow: hidden; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-lg); background: var(--courier-color-surface-raised); box-shadow: var(--courier-shadow); }
    .state-copy { display: grid; align-content: center; justify-items: start; gap: 1rem; padding: clamp(1.5rem, 5vw, 3rem); }
    .state-art { position: relative; min-height: 20rem; overflow: hidden; background: var(--courier-graphite-900); }
    .state-art courier-mascot { position: absolute; inset: 0; width: 100%; height: 100%; }
    .state-art courier-mascot::part(image) { width: 100%; height: 100%; object-fit: cover; }
    .empty, .loading { display: grid; min-height: 13rem; place-items: center; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-md); color: var(--courier-color-muted); background: var(--courier-color-surface-raised); font-family: var(--courier-font-mono); }
    @media (max-width: 64rem) { form { grid-template-columns: repeat(2, minmax(10rem, 1fr)); } }
    @media (max-width: 44rem) {
      header { align-items: flex-start; padding: 1rem 0; }
      nav { justify-content: flex-end; }
      .page-head, .server-head { grid-template-columns: 1fr; }
      .metrics { grid-template-columns: 1fr; }
      .metric + .metric { border-top: 1px solid var(--courier-color-border); border-left: 0; }
      form { grid-template-columns: 1fr; }
      dl { grid-template-columns: 1fr; }
      dt { margin-top: 0.3rem; }
      .state-brief { grid-template-columns: 1fr; }
      .state-art { min-height: 15rem; }
    }
  `]}connectedCallback(){super.connectedCallback(),this.theme=Ue(),this.unsubscribe=tt(e=>{this.snapshot=e,this.failed=!1}),this.refresh()}disconnectedCallback(){this.unsubscribe?.(),this.theme?.destroy(),super.disconnectedCallback()}async refresh(){this.failed=!1,this.conflict=!1;try{this.snapshot=await Qe()}catch{this.failed=!0}}setLocale(e){this.locale=e.detail}async stop(e,t){try{await et(e,t),await this.refresh()}catch{this.failed=!0}}async save(e,t){e.preventDefault();let n=new FormData(e.currentTarget),r={...t.policy,version:t.policy.version+1,auth:String(n.get(`auth`)),authAttempts:Number(n.get(`attempts`)),authFailAction:String(n.get(`failAction`)),noUi:n.get(`noUi`)===`on`};this.failed=!1,this.conflict=!1;try{await $e(t,r),await this.refresh()}catch(e){this.conflict=e instanceof Error&&e.name===`ConflictError`,this.failed=!this.conflict}}t(e){return rt(this.locale,e)}delivery(e){let t=e.source||this.t(`unavailable`),n=e.destination||this.t(`unavailable`);return L`
      <article>
        <div class="row delivery-head"><div><span class="label">${this.t(`deliveries`)} · ${e.route}</span><strong class="delivery-id">${e.id}</strong></div><courier-button @click=${()=>this.stop(`deliveries`,e.id)}>${this.t(`stopDelivery`)}</courier-button></div>
        <div class="route"><courier-route source=${t} destination=${n}></courier-route></div>
        <dl>
          <dt>${this.t(`source`)}</dt><dd>${t}</dd>
          <dt>${this.t(`destination`)}</dt><dd>${n}</dd>
          <dt>${this.t(`transferred`)}</dt><dd>${e.counters.confirmed}</dd>
        </dl>
        <span class="label">${this.t(`policy`)}</span>
        <form @submit=${t=>this.save(t,e)}>
          <label>${this.t(`authentication`)}<select name="auth"><option ?selected=${e.policy.auth===`none`}>none</option><option ?selected=${e.policy.auth===`basic`}>basic</option><option ?selected=${e.policy.auth===`password`}>password</option></select></label>
          <label>${this.t(`attempts`)}<input name="attempts" type="number" min="1" .value=${String(e.policy.authAttempts)}></label>
          <label>${this.t(`failAction`)}<select name="failAction"><option ?selected=${e.policy.authFailAction===`ban`}>ban</option><option ?selected=${e.policy.authFailAction===`stop`}>stop</option></select></label>
          <label class="checkbox"><input name="noUi" type="checkbox" ?checked=${e.policy.noUi}><span>${this.t(`noUi`)}</span></label>
          <courier-button type="submit" variant="primary">${this.t(`save`)}</courier-button>
        </form>
      </article>
    `}server(e){return L`
      <courier-panel>
        <section class="server">
          <div class="server-head">
            <div class="server-title"><span class="label">${this.t(`server`)}</span><h2 class="server-id">${e.id}</h2><span class="bind"><courier-icon name="server"></courier-icon>${e.bind}</span></div>
            <div class="actions"><courier-status tone=${e.status===`live`?`signal`:`danger`}>${this.t(e.status)}</courier-status><courier-button @click=${()=>this.stop(`servers`,e.id)}>${this.t(`stopServer`)}</courier-button></div>
          </div>
          <span class="label">${this.t(`deliveries`)}</span>
          <div class="delivery-stack">${e.deliveries.map(e=>this.delivery(e))}</div>
        </section>
      </courier-panel>
    `}render(){let e=this.snapshot?.servers??[],t=e.reduce((e,t)=>e+t.deliveries.length,0),n=e.reduce((e,t)=>e+t.deliveries.reduce((e,t)=>e+t.counters.confirmed,0),0);return L`
      <main>
        <header>
          <courier-brand product=${this.t(`brandProduct`)}></courier-brand>
          <nav><courier-button @click=${this.refresh}>${this.t(`refresh`)}</courier-button><courier-theme-selector .locale=${this.locale}></courier-theme-selector><courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector></nav>
        </header>
        <div class="workspace">
          <div class="page-head"><div><span class="eyebrow">${this.t(`eyebrow`)}</span><h1>${this.t(`title`)}</h1><p class="intro">${this.t(`intro`)}</p></div><courier-status tone=${this.failed?`danger`:`signal`}>${this.t(this.failed?`unreachable`:`live`)}</courier-status></div>
          <div class="metrics">
            <div class="metric"><strong>${e.length}</strong><span>${this.t(`serversMetric`)}</span></div>
            <div class="metric"><strong>${t}</strong><span>${this.t(`deliveriesMetric`)}</span></div>
            <div class="metric"><strong>${n}</strong><span>${this.t(`confirmedMetric`)}</span></div>
          </div>
          ${this.failed?L`<div class="state-brief"><div class="state-copy"><span class="eyebrow">${this.t(`unreachable`)}</span><p role="alert">${this.t(`failed`)}</p><courier-button @click=${this.refresh}>${this.t(`retry`)}</courier-button></div><div class="state-art"><courier-mascot alt="" .source=${Ze}></courier-mascot></div></div>`:z}
          ${this.conflict?L`<div class="notice"><p role="alert">${this.t(`conflict`)}</p><courier-button @click=${this.refresh}>${this.t(`refresh`)}</courier-button></div>`:z}
          ${this.snapshot?e.length===0?L`<div class="state-brief"><div class="state-copy"><span class="eyebrow">${this.t(`live`)}</span><p>${this.t(`empty`)}</p></div><div class="state-art"><courier-mascot alt="" .source=${Ze}></courier-mascot></div></div>`:L`<div class="server-list">${e.map(e=>this.server(e))}</div>`:this.failed?z:L`<div class="loading"><courier-status>${this.t(`loading`)}</courier-status></div>`}
        </div>
      </main>
    `}};customElements.get(`courier-admin-app`)||customElements.define(`courier-admin-app`,it);