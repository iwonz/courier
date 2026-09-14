(function(){let e=document.createElement(`link`).relList;if(e&&e.supports&&e.supports(`modulepreload`))return;for(let e of document.querySelectorAll(`link[rel="modulepreload"]`))n(e);new MutationObserver(e=>{for(let t of e)if(t.type===`childList`)for(let e of t.addedNodes)e.tagName===`LINK`&&e.rel===`modulepreload`&&n(e)}).observe(document,{childList:!0,subtree:!0});function t(e){let t={};return e.integrity&&(t.integrity=e.integrity),e.referrerPolicy&&(t.referrerPolicy=e.referrerPolicy),t.credentials=e.crossOrigin===`use-credentials`?`include`:e.crossOrigin===`anonymous`?`omit`:`same-origin`,t}function n(e){if(e.ep)return;e.ep=!0;let n=t(e);fetch(e.href,n)}})();var e=globalThis,t=e.ShadowRoot&&(e.ShadyCSS===void 0||e.ShadyCSS.nativeShadow)&&`adoptedStyleSheets`in Document.prototype&&`replace`in CSSStyleSheet.prototype,n=Symbol(),r=new WeakMap,i=class{constructor(e,t,r){if(this._$cssResult$=!0,r!==n)throw Error("CSSResult is not constructable. Use `unsafeCSS` or `css` instead.");this.cssText=e,this.t=t}get styleSheet(){let e=this.o,n=this.t;if(t&&e===void 0){let t=n!==void 0&&n.length===1;t&&(e=r.get(n)),e===void 0&&((this.o=e=new CSSStyleSheet).replaceSync(this.cssText),t&&r.set(n,e))}return e}toString(){return this.cssText}},a=e=>new i(typeof e==`string`?e:e+``,void 0,n),o=(e,...t)=>new i(e.length===1?e[0]:t.reduce((t,n,r)=>t+(e=>{if(!0===e._$cssResult$)return e.cssText;if(typeof e==`number`)return e;throw Error(`Value passed to 'css' function must be a 'css' function result: `+e+`. Use 'unsafeCSS' to pass non-literal values, but take care to ensure page security.`)})(n)+e[r+1],e[0]),e,n),s=(n,r)=>{if(t)n.adoptedStyleSheets=r.map(e=>e instanceof CSSStyleSheet?e:e.styleSheet);else for(let t of r){let r=document.createElement(`style`),i=e.litNonce;i!==void 0&&r.setAttribute(`nonce`,i),r.textContent=t.cssText,n.appendChild(r)}},c=t?e=>e:e=>e instanceof CSSStyleSheet?(e=>{let t=``;for(let n of e.cssRules)t+=n.cssText;return a(t)})(e):e,{is:l,defineProperty:u,getOwnPropertyDescriptor:d,getOwnPropertyNames:ee,getOwnPropertySymbols:te,getPrototypeOf:ne}=Object,f=globalThis,re=f.trustedTypes,ie=re?re.emptyScript:``,ae=f.reactiveElementPolyfillSupport,p=(e,t)=>e,m={toAttribute(e,t){switch(t){case Boolean:e=e?ie:null;break;case Object:case Array:e=e==null?e:JSON.stringify(e)}return e},fromAttribute(e,t){let n=e;switch(t){case Boolean:n=e!==null;break;case Number:n=e===null?null:Number(e);break;case Object:case Array:try{n=JSON.parse(e)}catch{n=null}}return n}},h=(e,t)=>!l(e,t),g={attribute:!0,type:String,converter:m,reflect:!1,useDefault:!1,hasChanged:h};Symbol.metadata??=Symbol(`metadata`),f.litPropertyMetadata??=new WeakMap;var _=class extends HTMLElement{static addInitializer(e){this._$Ei(),(this.l??=[]).push(e)}static get observedAttributes(){return this.finalize(),this._$Eh&&[...this._$Eh.keys()]}static createProperty(e,t=g){if(t.state&&(t.attribute=!1),this._$Ei(),this.prototype.hasOwnProperty(e)&&((t=Object.create(t)).wrapped=!0),this.elementProperties.set(e,t),!t.noAccessor){let n=Symbol(),r=this.getPropertyDescriptor(e,n,t);r!==void 0&&u(this.prototype,e,r)}}static getPropertyDescriptor(e,t,n){let{get:r,set:i}=d(this.prototype,e)??{get(){return this[t]},set(e){this[t]=e}};return{get:r,set(t){let a=r?.call(this);i?.call(this,t),this.requestUpdate(e,a,n)},configurable:!0,enumerable:!0}}static getPropertyOptions(e){return this.elementProperties.get(e)??g}static _$Ei(){if(this.hasOwnProperty(p(`elementProperties`)))return;let e=ne(this);e.finalize(),e.l!==void 0&&(this.l=[...e.l]),this.elementProperties=new Map(e.elementProperties)}static finalize(){if(this.hasOwnProperty(p(`finalized`)))return;if(this.finalized=!0,this._$Ei(),this.hasOwnProperty(p(`properties`))){let e=this.properties,t=[...ee(e),...te(e)];for(let n of t)this.createProperty(n,e[n])}let e=this[Symbol.metadata];if(e!==null){let t=litPropertyMetadata.get(e);if(t!==void 0)for(let[e,n]of t)this.elementProperties.set(e,n)}this._$Eh=new Map;for(let[e,t]of this.elementProperties){let n=this._$Eu(e,t);n!==void 0&&this._$Eh.set(n,e)}this.elementStyles=this.finalizeStyles(this.styles)}static finalizeStyles(e){let t=[];if(Array.isArray(e)){let n=new Set(e.flat(1/0).reverse());for(let e of n)t.unshift(c(e))}else e!==void 0&&t.push(c(e));return t}static _$Eu(e,t){let n=t.attribute;return!1===n?void 0:typeof n==`string`?n:typeof e==`string`?e.toLowerCase():void 0}constructor(){super(),this._$Ep=void 0,this.isUpdatePending=!1,this.hasUpdated=!1,this._$Em=null,this._$Ev()}_$Ev(){this._$ES=new Promise(e=>this.enableUpdating=e),this._$AL=new Map,this._$E_(),this.requestUpdate(),this.constructor.l?.forEach(e=>e(this))}addController(e){(this._$EO??=new Set).add(e),this.renderRoot!==void 0&&this.isConnected&&e.hostConnected?.()}removeController(e){this._$EO?.delete(e)}_$E_(){let e=new Map,t=this.constructor.elementProperties;for(let n of t.keys())this.hasOwnProperty(n)&&(e.set(n,this[n]),delete this[n]);e.size>0&&(this._$Ep=e)}createRenderRoot(){let e=this.shadowRoot??this.attachShadow(this.constructor.shadowRootOptions);return s(e,this.constructor.elementStyles),e}connectedCallback(){this.renderRoot??=this.createRenderRoot(),this.enableUpdating(!0),this._$EO?.forEach(e=>e.hostConnected?.())}enableUpdating(e){}disconnectedCallback(){this._$EO?.forEach(e=>e.hostDisconnected?.())}attributeChangedCallback(e,t,n){this._$AK(e,n)}_$ET(e,t){let n=this.constructor.elementProperties.get(e),r=this.constructor._$Eu(e,n);if(r!==void 0&&!0===n.reflect){let i=(n.converter?.toAttribute===void 0?m:n.converter).toAttribute(t,n.type);this._$Em=e,i==null?this.removeAttribute(r):this.setAttribute(r,i),this._$Em=null}}_$AK(e,t){let n=this.constructor,r=n._$Eh.get(e);if(r!==void 0&&this._$Em!==r){let e=n.getPropertyOptions(r),i=typeof e.converter==`function`?{fromAttribute:e.converter}:e.converter?.fromAttribute===void 0?m:e.converter;this._$Em=r;let a=i.fromAttribute(t,e.type);this[r]=a??this._$Ej?.get(r)??a,this._$Em=null}}requestUpdate(e,t,n,r=!1,i){if(e!==void 0){let a=this.constructor;if(!1===r&&(i=this[e]),n??=a.getPropertyOptions(e),!((n.hasChanged??h)(i,t)||n.useDefault&&n.reflect&&i===this._$Ej?.get(e)&&!this.hasAttribute(a._$Eu(e,n))))return;this.C(e,t,n)}!1===this.isUpdatePending&&(this._$ES=this._$EP())}C(e,t,{useDefault:n,reflect:r,wrapped:i},a){n&&!(this._$Ej??=new Map).has(e)&&(this._$Ej.set(e,a??t??this[e]),!0!==i||a!==void 0)||(this._$AL.has(e)||(this.hasUpdated||n||(t=void 0),this._$AL.set(e,t)),!0===r&&this._$Em!==e&&(this._$Eq??=new Set).add(e))}async _$EP(){this.isUpdatePending=!0;try{await this._$ES}catch(e){Promise.reject(e)}let e=this.scheduleUpdate();return e!=null&&await e,!this.isUpdatePending}scheduleUpdate(){return this.performUpdate()}performUpdate(){if(!this.isUpdatePending)return;if(!this.hasUpdated){if(this.renderRoot??=this.createRenderRoot(),this._$Ep){for(let[e,t]of this._$Ep)this[e]=t;this._$Ep=void 0}let e=this.constructor.elementProperties;if(e.size>0)for(let[t,n]of e){let{wrapped:e}=n,r=this[t];!0!==e||this._$AL.has(t)||r===void 0||this.C(t,void 0,n,r)}}let e=!1,t=this._$AL;try{e=this.shouldUpdate(t),e?(this.willUpdate(t),this._$EO?.forEach(e=>e.hostUpdate?.()),this.update(t)):this._$EM()}catch(t){throw e=!1,this._$EM(),t}e&&this._$AE(t)}willUpdate(e){}_$AE(e){this._$EO?.forEach(e=>e.hostUpdated?.()),this.hasUpdated||(this.hasUpdated=!0,this.firstUpdated(e)),this.updated(e)}_$EM(){this._$AL=new Map,this.isUpdatePending=!1}get updateComplete(){return this.getUpdateComplete()}getUpdateComplete(){return this._$ES}shouldUpdate(e){return!0}update(e){this._$Eq&&=this._$Eq.forEach(e=>this._$ET(e,this[e])),this._$EM()}updated(e){}firstUpdated(e){}};_.elementStyles=[],_.shadowRootOptions={mode:`open`},_[p(`elementProperties`)]=new Map,_[p(`finalized`)]=new Map,ae?.({ReactiveElement:_}),(f.reactiveElementVersions??=[]).push(`2.1.2`);var v=globalThis,y=e=>e,b=v.trustedTypes,x=b?b.createPolicy(`lit-html`,{createHTML:e=>e}):void 0,S=`$lit$`,C=`lit$${Math.random().toFixed(9).slice(2)}$`,oe=`?`+C,se=`<${oe}>`,w=document,T=()=>w.createComment(``),E=e=>e===null||typeof e!=`object`&&typeof e!=`function`,D=Array.isArray,ce=e=>D(e)||typeof e?.[Symbol.iterator]==`function`,O=`[ 	
\f\r]`,k=/<(?:(!--|\/[^a-zA-Z])|(\/?[a-zA-Z][^>\s]*)|(\/?$))/g,A=/-->/g,j=/>/g,M=RegExp(`>|${O}(?:([^\\s"'>=/]+)(${O}*=${O}*(?:[^ \t\n\f\r"'\`<>=]|("|')|))|$)`,`g`),N=/'/g,le=/"/g,P=/^(?:script|style|textarea|title)$/i,F=(e=>(t,...n)=>({_$litType$:e,strings:t,values:n}))(1),I=Symbol.for(`lit-noChange`),L=Symbol.for(`lit-nothing`),R=new WeakMap,z=w.createTreeWalker(w,129);function B(e,t){if(!D(e)||!e.hasOwnProperty(`raw`))throw Error(`invalid template strings array`);return x===void 0?t:x.createHTML(t)}var ue=(e,t)=>{let n=e.length-1,r=[],i,a=t===2?`<svg>`:t===3?`<math>`:``,o=k;for(let t=0;t<n;t++){let n=e[t],s,c,l=-1,u=0;for(;u<n.length&&(o.lastIndex=u,c=o.exec(n),c!==null);)u=o.lastIndex,o===k?c[1]===`!--`?o=A:c[1]===void 0?c[2]===void 0?c[3]!==void 0&&(o=M):(P.test(c[2])&&(i=RegExp(`</`+c[2],`g`)),o=M):o=j:o===M?c[0]===`>`?(o=i??k,l=-1):c[1]===void 0?l=-2:(l=o.lastIndex-c[2].length,s=c[1],o=c[3]===void 0?M:c[3]===`"`?le:N):o===le||o===N?o=M:o===A||o===j?o=k:(o=M,i=void 0);let d=o===M&&e[t+1].startsWith(`/>`)?` `:``;a+=o===k?n+se:l>=0?(r.push(s),n.slice(0,l)+S+n.slice(l)+C+d):n+C+(l===-2?t:d)}return[B(e,a+(e[n]||`<?>`)+(t===2?`</svg>`:t===3?`</math>`:``)),r]},V=class e{constructor({strings:t,_$litType$:n},r){let i;this.parts=[];let a=0,o=0,s=t.length-1,c=this.parts,[l,u]=ue(t,n);if(this.el=e.createElement(l,r),z.currentNode=this.el.content,n===2||n===3){let e=this.el.content.firstChild;e.replaceWith(...e.childNodes)}for(;(i=z.nextNode())!==null&&c.length<s;){if(i.nodeType===1){if(i.hasAttributes())for(let e of i.getAttributeNames())if(e.endsWith(S)){let t=u[o++],n=i.getAttribute(e).split(C),r=/([.?@])?(.*)/.exec(t);c.push({type:1,index:a,name:r[2],strings:n,ctor:r[1]===`.`?fe:r[1]===`?`?pe:r[1]===`@`?me:W}),i.removeAttribute(e)}else e.startsWith(C)&&(c.push({type:6,index:a}),i.removeAttribute(e));if(P.test(i.tagName)){let e=i.textContent.split(C),t=e.length-1;if(t>0){i.textContent=b?b.emptyScript:``;for(let n=0;n<t;n++)i.append(e[n],T()),z.nextNode(),c.push({type:2,index:++a});i.append(e[t],T())}}}else if(i.nodeType===8){if(i.data===oe)c.push({type:2,index:a});else{let e=-1;for(;(e=i.data.indexOf(C,e+1))!==-1;)c.push({type:7,index:a}),e+=C.length-1}}a++}}static createElement(e,t){let n=w.createElement(`template`);return n.innerHTML=e,n}};function H(e,t,n=e,r){if(t===I)return t;let i=r===void 0?n._$Cl:n._$Co?.[r],a=E(t)?void 0:t._$litDirective$;return i?.constructor!==a&&(i?._$AO?.(!1),a===void 0?i=void 0:(i=new a(e),i._$AT(e,n,r)),r===void 0?n._$Cl=i:(n._$Co??=[])[r]=i),i!==void 0&&(t=H(e,i._$AS(e,t.values),i,r)),t}var de=class{constructor(e,t){this._$AV=[],this._$AN=void 0,this._$AD=e,this._$AM=t}get parentNode(){return this._$AM.parentNode}get _$AU(){return this._$AM._$AU}u(e){let{el:{content:t},parts:n}=this._$AD,r=(e?.creationScope??w).importNode(t,!0);z.currentNode=r;let i=z.nextNode(),a=0,o=0,s=n[0];for(;s!==void 0;){if(a===s.index){let t;s.type===2?t=new U(i,i.nextSibling,this,e):s.type===1?t=new s.ctor(i,s.name,s.strings,this,e):s.type===6&&(t=new he(i,this,e)),this._$AV.push(t),s=n[++o]}a!==s?.index&&(i=z.nextNode(),a++)}return z.currentNode=w,r}p(e){let t=0;for(let n of this._$AV)n!==void 0&&(n.strings===void 0?n._$AI(e[t]):(n._$AI(e,n,t),t+=n.strings.length-2)),t++}},U=class e{get _$AU(){return this._$AM?._$AU??this._$Cv}constructor(e,t,n,r){this.type=2,this._$AH=L,this._$AN=void 0,this._$AA=e,this._$AB=t,this._$AM=n,this.options=r,this._$Cv=r?.isConnected??!0}get parentNode(){let e=this._$AA.parentNode,t=this._$AM;return t!==void 0&&e?.nodeType===11&&(e=t.parentNode),e}get startNode(){return this._$AA}get endNode(){return this._$AB}_$AI(e,t=this){e=H(this,e,t),E(e)?e===L||e==null||e===``?(this._$AH!==L&&this._$AR(),this._$AH=L):e!==this._$AH&&e!==I&&this._(e):e._$litType$===void 0?e.nodeType===void 0?ce(e)?this.k(e):this._(e):this.T(e):this.$(e)}O(e){return this._$AA.parentNode.insertBefore(e,this._$AB)}T(e){this._$AH!==e&&(this._$AR(),this._$AH=this.O(e))}_(e){this._$AH!==L&&E(this._$AH)?this._$AA.nextSibling.data=e:this.T(w.createTextNode(e)),this._$AH=e}$(e){let{values:t,_$litType$:n}=e,r=typeof n==`number`?this._$AC(e):(n.el===void 0&&(n.el=V.createElement(B(n.h,n.h[0]),this.options)),n);if(this._$AH?._$AD===r)this._$AH.p(t);else{let e=new de(r,this),n=e.u(this.options);e.p(t),this.T(n),this._$AH=e}}_$AC(e){let t=R.get(e.strings);return t===void 0&&R.set(e.strings,t=new V(e)),t}k(t){D(this._$AH)||(this._$AH=[],this._$AR());let n=this._$AH,r,i=0;for(let a of t)i===n.length?n.push(r=new e(this.O(T()),this.O(T()),this,this.options)):r=n[i],r._$AI(a),i++;i<n.length&&(this._$AR(r&&r._$AB.nextSibling,i),n.length=i)}_$AR(e=this._$AA.nextSibling,t){for(this._$AP?.(!1,!0,t);e!==this._$AB;){let t=y(e).nextSibling;y(e).remove(),e=t}}setConnected(e){this._$AM===void 0&&(this._$Cv=e,this._$AP?.(e))}},W=class{get tagName(){return this.element.tagName}get _$AU(){return this._$AM._$AU}constructor(e,t,n,r,i){this.type=1,this._$AH=L,this._$AN=void 0,this.element=e,this.name=t,this._$AM=r,this.options=i,n.length>2||n[0]!==``||n[1]!==``?(this._$AH=Array(n.length-1).fill(new String),this.strings=n):this._$AH=L}_$AI(e,t=this,n,r){let i=this.strings,a=!1;if(i===void 0)e=H(this,e,t,0),a=!E(e)||e!==this._$AH&&e!==I,a&&(this._$AH=e);else{let r=e,o,s;for(e=i[0],o=0;o<i.length-1;o++)s=H(this,r[n+o],t,o),s===I&&(s=this._$AH[o]),a||=!E(s)||s!==this._$AH[o],s===L?e=L:e!==L&&(e+=(s??``)+i[o+1]),this._$AH[o]=s}a&&!r&&this.j(e)}j(e){e===L?this.element.removeAttribute(this.name):this.element.setAttribute(this.name,e??``)}},fe=class extends W{constructor(){super(...arguments),this.type=3}j(e){this.element[this.name]=e===L?void 0:e}},pe=class extends W{constructor(){super(...arguments),this.type=4}j(e){this.element.toggleAttribute(this.name,!!e&&e!==L)}},me=class extends W{constructor(e,t,n,r,i){super(e,t,n,r,i),this.type=5}_$AI(e,t=this){if((e=H(this,e,t,0)??L)===I)return;let n=this._$AH,r=e===L&&n!==L||e.capture!==n.capture||e.once!==n.once||e.passive!==n.passive,i=e!==L&&(n===L||r);r&&this.element.removeEventListener(this.name,this,n),i&&this.element.addEventListener(this.name,this,e),this._$AH=e}handleEvent(e){typeof this._$AH==`function`?this._$AH.call(this.options?.host??this.element,e):this._$AH.handleEvent(e)}},he=class{constructor(e,t,n){this.element=e,this.type=6,this._$AN=void 0,this._$AM=t,this.options=n}get _$AU(){return this._$AM._$AU}_$AI(e){H(this,e)}},ge=v.litHtmlPolyfillSupport;ge?.(V,U),(v.litHtmlVersions??=[]).push(`3.3.3`);var _e=(e,t,n)=>{let r=n?.renderBefore??t,i=r._$litPart$;if(i===void 0){let e=n?.renderBefore??null;r._$litPart$=i=new U(t.insertBefore(T(),e),e,void 0,n??{})}return i._$AI(e),i},G=globalThis,K=class extends _{constructor(){super(...arguments),this.renderOptions={host:this},this._$Do=void 0}createRenderRoot(){let e=super.createRenderRoot();return this.renderOptions.renderBefore??=e.firstChild,e}update(e){let t=this.render();this.hasUpdated||(this.renderOptions.isConnected=this.isConnected),super.update(e),this._$Do=_e(t,this.renderRoot,this.renderOptions)}connectedCallback(){super.connectedCallback(),this._$Do?.setConnected(!0)}disconnectedCallback(){super.disconnectedCallback(),this._$Do?.setConnected(!1)}render(){return I}};K._$litElement$=!0,K.finalized=!0,G.litElementHydrateSupport?.({LitElement:K});var ve=G.litElementPolyfillSupport;ve?.({LitElement:K}),(G.litElementVersions??=[]).push(`4.2.2`);var ye=o`
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
`;var be=o`
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
`,xe=class extends K{constructor(...e){super(...e),this.disabled=!1,this.type=`button`,this.variant=`secondary`}static{this.properties={disabled:{type:Boolean,reflect:!0},type:{type:String,reflect:!0},variant:{type:String,reflect:!0}}}static{this.styles=[ye,o`
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
  `]}render(){return F`<button type=${this.type} ?disabled=${this.disabled}><slot></slot></button>`}},Se=class extends K{constructor(...e){super(...e),this.product=``}static{this.properties={product:{type:String}}}static{this.styles=o`
    :host { display: inline-flex; min-width: 0; color: var(--courier-color-text, #151714); font-family: var(--courier-font-sans, sans-serif); }
    .lockup { display: inline-flex; min-width: 0; align-items: center; gap: 0.625rem; color: inherit; }
    svg { flex: 0 0 auto; width: 2.25rem; height: 2.25rem; }
    .words { display: grid; line-height: 1; }
    strong { font-family: var(--courier-font-display, sans-serif); font-size: 1.125rem; font-weight: 850; letter-spacing: -0.04em; }
    small { margin-top: 0.25rem; color: var(--courier-color-muted, #596054); font-family: var(--courier-font-mono, monospace); font-size: 0.625rem; font-weight: 700; letter-spacing: 0.09em; text-transform: uppercase; }
  `}render(){return F`<span class="lockup">
      <svg viewBox="0 0 40 40" aria-hidden="true">
        <rect x="1" y="1" width="38" height="38" rx="4" fill="var(--courier-color-inverse, #151714)"></rect>
        <path d="M10 9v22h8v-4h-4V13h4V9zm8 7h8v8h-8zm7-5 8 9-8 9v-6h-1v-6h1z" fill="var(--courier-color-accent, #d4ff45)"></path>
      </svg>
      <span class="words"><strong>Courier</strong><small>${this.product}</small></span>
    </span>`}},Ce=class extends K{constructor(...e){super(...e),this.alt=``,this.eager=!1,this.source=``}static{this.properties={alt:{type:String},eager:{type:Boolean},source:{type:String}}}static{this.styles=o`
    :host { display: block; }
    img { display: block; width: 100%; height: auto; filter: drop-shadow(0 1.5rem 2rem rgb(16 18 15 / 0.18)); }
  `}render(){return F`<img part="image" src=${this.source} alt=${this.alt} decoding="async" loading=${this.eager?`eager`:`lazy`} fetchpriority=${this.eager?`high`:`auto`}>`}},we=class extends K{constructor(...e){super(...e),this.tone=`neutral`}static{this.properties={tone:{type:String,reflect:!0}}}static{this.styles=o`
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
  `}render(){return F`<i aria-hidden="true"></i><slot></slot>`}},Te=class extends K{constructor(...e){super(...e),this.source=``,this.destination=``}static{this.properties={source:{type:String},destination:{type:String}}}static{this.styles=o`
    :host { display: grid; color: var(--courier-color-text, #151714); font-family: var(--courier-font-mono, monospace); }
    .route { display: grid; grid-template-columns: minmax(0, 1fr) minmax(3rem, 0.55fr) minmax(0, 1fr); align-items: center; gap: 0.65rem; }
    .node { overflow: hidden; padding: 0.65rem 0.75rem; border: 1px solid var(--courier-color-border, #c8cdbf); border-radius: var(--courier-radius-sm, 0.25rem); background: var(--courier-color-surface, #fafbf3); font-size: 0.75rem; text-overflow: ellipsis; white-space: nowrap; }
    .line { position: relative; height: 1px; color: var(--courier-color-border-strong, #8e9587); background: currentColor; }
    .line::before { content: ""; position: absolute; top: -0.2rem; left: 0; width: 0.45rem; height: 0.45rem; border-radius: 50%; background: var(--courier-color-accent, #d4ff45); }
    .line::after { content: ""; position: absolute; top: -0.22rem; right: 0; width: 0.4rem; height: 0.4rem; border-top: 1px solid currentColor; border-right: 1px solid currentColor; transform: rotate(45deg); }
  `}render(){return F`<div class="route"><span class="node">${this.source}</span><span class="line" aria-hidden="true"></span><span class="node">${this.destination}</span></div>`}},q={en:{"theme.label":`Theme`,"theme.system":`System`,"theme.light":`Light`,"theme.dark":`Dark`,"locale.label":`Language`,"locale.en":`English`,"locale.ru":`Russian`,"progress.label":`Delivery progress`,"action.cancel":`Cancel`,"action.close":`Close`},ru:{"theme.label":`Тема`,"theme.system":`Системная`,"theme.light":`Светлая`,"theme.dark":`Тёмная`,"locale.label":`Язык`,"locale.en":`Английский`,"locale.ru":`Русский`,"progress.label":`Ход доставки`,"action.cancel":`Отмена`,"action.close":`Закрыть`}},Ee=Object.keys(q),J=`courier.locale`;function Y(e){if(!e)return;let t=e.toLowerCase().split(`-`)[0];return Ee.includes(t)?t:void 0}function De(e){for(let t of e){let e=Y(t);if(e)return e}return`en`}function Oe(e,t){if(e)try{let t=Y(e.getItem(J));if(t)return t}catch{}return De(t)}function ke(e,t){if(e)try{e.setItem(J,t)}catch{}}function X(e,t){return q[Y(e)??`en`][t]}function Ae(){let e;try{e=globalThis.localStorage}catch{e=void 0}let t=globalThis.navigator?.languages??[];return Oe(e,t)}var je=class extends K{constructor(...e){super(...e),this.locale=`en`}static{this.properties={locale:{type:String}}}connectedCallback(){super.connectedCallback(),this.locale=Ae()}change(e){let t=Y(e.detail)??`en`;this.locale=t;let n;try{n=globalThis.localStorage}catch{n=void 0}ke(n,t),this.dispatchEvent(new CustomEvent(`courier-locale-change`,{detail:t,bubbles:!0,composed:!0}))}render(){return F`<courier-segmented-control
      .label=${X(this.locale,`locale.label`)}
      .value=${this.locale}
      .options=${[{value:`en`,label:X(this.locale,`locale.en`)},{value:`ru`,label:X(this.locale,`locale.ru`)}]}
      data-storage-key=${J}
      @courier-segment-change=${this.change}
    ></courier-segmented-control>`}},Me=class extends K{constructor(...e){super(...e),this.heading=``}static{this.properties={heading:{type:String}}}static{this.styles=o`
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
  `}render(){return this.heading?F`<section aria-labelledby="courier-panel-heading">
      <h2 id="courier-panel-heading">${this.heading}</h2>
      <slot></slot>
    </section>`:F`<section><slot></slot></section>`}};function Ne(e,t){return!Number.isFinite(e)||!Number.isFinite(t)||t<=0?0:Math.min(1,Math.max(0,e/t))}var Pe=class extends K{constructor(...e){super(...e),this.value=0,this.total=0,this.label=``,this.locale=`en`}static{this.properties={value:{type:Number},total:{type:Number},label:{type:String},locale:{type:String}}}static{this.styles=o`
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
  `}render(){let e=Ne(this.value,this.total),t=this.label||X(this.locale,`progress.label`),n=Number.isFinite(this.total)&&this.total>0?this.total:0;return F`<div
      class="track"
      role="progressbar"
      aria-label=${t}
      aria-valuemin="0"
      aria-valuemax=${n}
      aria-valuenow=${Number.isFinite(this.value)?Math.max(0,Math.min(this.value,n)):0}
    ><div class="fill" style=${`transform: scaleX(${e})`}></div></div>
    <output>${Math.round(e*100)}%</output>`}};function Fe(e,t,n){if(!(n<=0))switch(e){case`ArrowLeft`:case`ArrowUp`:return(t-1+n)%n;case`ArrowRight`:case`ArrowDown`:return(t+1)%n;case`Home`:return 0;case`End`:return n-1;default:return}}var Ie=class extends K{constructor(...e){super(...e),this.label=``,this.options=[],this.value=``}static{this.properties={label:{type:String},options:{attribute:!1},value:{type:String}}}static{this.styles=o`
    :host { display: block; min-width: 0; color: var(--courier-color-text, #151714); font-family: var(--courier-font-sans, sans-serif); }
    fieldset { min-width: 0; margin: 0; padding: 0; border: 0; }
    legend { margin: 0 0 0.25rem; padding: 0; color: var(--courier-color-muted, #596054); font-family: var(--courier-font-mono, monospace); font-size: 0.625rem; font-weight: 750; letter-spacing: 0.08em; line-height: 1; text-transform: uppercase; }
    .segments { display: inline-grid; max-width: 100%; grid-auto-columns: minmax(0, auto); grid-auto-flow: column; gap: 2px; padding: 2px; border: 1px solid var(--courier-color-border, #c8cdbf); border-radius: var(--courier-radius-md, 0.625rem); background: color-mix(in srgb, var(--courier-color-field, #e7e9dc) 68%, transparent); box-shadow: inset 0 1px 2px rgb(16 18 15 / 0.07); }
    button { appearance: none; min-width: 0; min-height: 2.25rem; padding: 0.45rem 0.68rem; overflow: hidden; border: 1px solid transparent; border-radius: calc(var(--courier-radius-md, 0.625rem) - 3px); color: var(--courier-color-muted, #596054); background: transparent; font: inherit; font-size: 0.75rem; font-weight: 780; line-height: 1; text-overflow: ellipsis; white-space: nowrap; cursor: pointer; transition: color var(--courier-duration, 160ms) var(--courier-ease, ease), background var(--courier-duration, 160ms) var(--courier-ease, ease), transform var(--courier-duration, 160ms) var(--courier-ease, ease); }
    button:hover { color: var(--courier-color-text, #151714); background: color-mix(in srgb, var(--courier-color-surface-raised, #fff) 72%, transparent); }
    button.selected { border-color: color-mix(in srgb, var(--courier-color-accent, #d4ff45) 64%, var(--courier-color-border, #c8cdbf)); color: var(--courier-color-accent-ink, #151714); background: var(--courier-color-accent, #d4ff45); box-shadow: 0 1px 0 rgb(16 18 15 / 0.12); }
    button:active { transform: translateY(1px); }
    button:focus-visible { position: relative; z-index: 1; outline: 3px solid var(--courier-beak, #ff8758); outline-offset: 2px; }
  `}activate(e){let t=e.currentTarget.dataset.value??``;t&&t!==this.value&&(this.value=t,this.dispatchEvent(new CustomEvent(`courier-segment-change`,{detail:t,bubbles:!0,composed:!0})))}move(e){let t=[...this.renderRoot.querySelectorAll(`button`)],n=t.indexOf(e.currentTarget),r=Fe(e.key,n,t.length);if(r===void 0)return;e.preventDefault();let i=t[r];i.focus(),i.click()}render(){return F`<fieldset>
      <legend>${this.label}</legend>
      <div class="segments" role="radiogroup" aria-label=${this.label}>
        ${this.options.map(e=>{let t=e.value===this.value;return F`<button
            type="button"
            role="radio"
            class=${t?`selected`:``}
            data-value=${e.value}
            aria-checked=${String(t)}
            tabindex=${t?0:-1}
            @click=${this.activate}
            @keydown=${this.move}
          >${e.label}</button>`})}
      </div>
    </fieldset>`}},Le=[`system`,`light`,`dark`],Re=`courier.theme`;function ze(e){return Le.includes(e)?e:`system`}function Be(e,t){return e===`system`?t?.matches?`dark`:`light`:e}function Ve(e){if(!e)return`system`;try{return ze(e.getItem(Re))}catch{return`system`}}function He(e,t){if(e)try{e.setItem(Re,t)}catch{}}var Ue=class{constructor(e,t,n,r){this.root=e,this.storage=t,this.media=n,this.onSystemChange=()=>this.apply(),this.preference=r??Ve(t),this.media?.addEventListener(`change`,this.onSystemChange),this.apply()}set(e){this.preference=e,He(this.storage,e),this.apply()}destroy(){this.media?.removeEventListener(`change`,this.onSystemChange)}apply(){this.root.dataset.courierTheme=Be(this.preference,this.media),this.root.dataset.courierThemePreference=this.preference}};function We(){let e;try{e=globalThis.localStorage}catch{e=void 0}let t=globalThis.matchMedia?.(`(prefers-color-scheme: dark)`);return new Ue(document.documentElement,e,t)}var Ge=class extends K{constructor(...e){super(...e),this.preference=`system`,this.locale=`en`}static{this.properties={preference:{type:String},locale:{type:String}}}connectedCallback(){super.connectedCallback(),this.state=We(),this.preference=this.state.preference}disconnectedCallback(){this.state?.destroy(),super.disconnectedCallback()}change(e){let t=ze(e.detail);this.preference=t,this.state?.set(t),this.dispatchEvent(new CustomEvent(`courier-theme-change`,{detail:t,bubbles:!0,composed:!0}))}render(){return F`<courier-segmented-control
      .label=${X(this.locale,`theme.label`)}
      .value=${this.preference}
      .options=${[{value:`system`,label:X(this.locale,`theme.system`)},{value:`light`,label:X(this.locale,`theme.light`)},{value:`dark`,label:X(this.locale,`theme.dark`)}]}
      @courier-segment-change=${this.change}
    ></courier-segmented-control>`}},Ke=[`archive`,`copy`,`download`,`folder`,`parcel`,`receipt`,`retry`,`route`,`server`,`shield`,`upload`],qe={archive:`M3 3h18v5H3zM5 8v13h14V8M9 12h6`,copy:`M8 3h13v13M3 8h13v13H3z`,download:`M12 3v13m-5-5 5 5 5-5M4 15v6h16v-6`,folder:`M3 6h7l2 3h9v12H3zM3 6V3h7l2 3h9v3`,parcel:`m3 7 9-5 9 5v11l-9 4-9-4zM3 7l9 5 9-5M12 12v10M8 4l9 5v5`,receipt:`M5 2h14v20l-3-2-4 2-4-2-3 2zM8 7h8M8 11h8m-8 5 2 2 5-4`,retry:`M3 10a9 9 0 1 1 1 7M3 3v7h7M12 7v5l3 2`,route:`M2 3h6v6H2zM16 15h6v6h-6zM11 6h8v6m-3-3 3 3 3-3M13 18H5v-6m-3 3 3-3 3 3`,server:`M3 2h18v8H3zM3 14h18v8H3zM7 6h1m3 0h6M7 18h1m3 0h6M6 10v4m12-4v4`,shield:`m12 2 8 3v7c0 5-8 10-8 10S4 17 4 12V5zM8 11l3 3 5-6`,upload:`M12 16V3m-5 5 5-5 5 5M4 15v6h16v-6`};function Je(e){return Ke.includes(e)?e:`parcel`}var Ye={"courier-brand":Se,"courier-button":xe,"courier-icon":class extends K{constructor(...e){super(...e),this.name=`parcel`,this.label=``}static{this.properties={name:{type:String},label:{type:String}}}static{this.styles=o`
    :host {
      display: inline-flex;
      width: 1.5rem;
      height: 1.5rem;
      color: currentColor;
    }
    svg { width: 100%; height: 100%; }
  `}render(){let e=Je(this.name);return F`<svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="1.7"
      stroke-linecap="square"
      stroke-linejoin="miter"
      role=${this.label?`img`:`presentation`}
      aria-hidden=${this.label?`false`:`true`}
      aria-label=${this.label||void 0}
    ><path d=${qe[e]}></path></svg>`}},"courier-locale-selector":je,"courier-mascot":Ce,"courier-panel":Me,"courier-progress":Pe,"courier-route":Te,"courier-segmented-control":Ie,"courier-status":we,"courier-theme-selector":Ge};function Xe(e=customElements){for(let[t,n]of Object.entries(Ye))e.get(t)||e.define(t,n)}var Ze=`/assets/data.webp`;function Z(e,t=``,n=globalThis.location.pathname){let r=n.endsWith(`/`)?n:`${n}/`,i=new URL(`api/v1/${e}`,globalThis.location.origin);return i.pathname=`${r}api/v1/${e}`,t&&i.searchParams.set(`path`,t),`${i.pathname}${i.search}`}async function Q(e){if(!e.ok)throw Error(`Courier request failed (${e.status})`);return e.json()}async function Qe(e=``,t=globalThis.fetch){return Q(await t(Z(`meta`,e),{credentials:`same-origin`}))}async function $e(e,t=globalThis.fetch){return Q(await t(Z(`session`),{method:`POST`,credentials:`same-origin`,headers:{"Content-Type":`application/json`},body:JSON.stringify({password:e})}))}async function et(e,t,n=globalThis.fetch){let r=new FormData;r.append(`file`,e),await Q(await n(Z(`upload`),{method:`POST`,credentials:`same-origin`,headers:{"X-Courier-CSRF":t,"X-Courier-File-Size":String(e.size)},body:r}))}function $(e,t=!1){let n=new URL(Z(`download`,e),globalThis.location.origin);return t&&n.searchParams.set(`archive`,`tar.gz`),`${n.pathname}${n.search}`}function tt(e,t){return e?`${e}/${t}`:t}function nt(e){let t=e.lastIndexOf(`/`);return t<0?``:e.slice(0,t)}var rt={en:{brandProduct:`Delivery terminal`,title:`Courier delivery`,privateRoute:`Private route`,loading:`Preparing the delivery route…`,retry:`Retry connection`,download:`Download file`,downloadArchive:`Download as archive`,downloadAll:`Download directory`,upload:`Choose a file`,uploadTitle:`Dispatch a file`,uploadHelp:`Select one file. Courier checks policy and reserves the final path before committing it.`,up:`Parent directory`,password:`Delivery password`,signIn:`Verify access`,accessTitle:`Identity check`,accessHelp:`This route is protected. Verify access to reveal its delivery metadata.`,empty:`No entries are available at this path.`,failed:`The route is unavailable or authorization is required. No delivery metadata was revealed.`,ready:`Route ready`,manifest:`Delivery manifest`,confirmed:`Verified handoff`,entryTypeFile:`File`,entryTypeDirectory:`Directory`,itemSize:`Bytes`},ru:{brandProduct:`Терминал доставки`,title:`Доставка Courier`,privateRoute:`Приватный маршрут`,loading:`Подготовка маршрута доставки…`,retry:`Повторить подключение`,download:`Скачать файл`,downloadArchive:`Скачать архивом`,downloadAll:`Скачать каталог`,upload:`Выбрать файл`,uploadTitle:`Отправить файл`,uploadHelp:`Выберите один файл. Courier проверит правила и зарезервирует конечный путь до фиксации.`,up:`Родительский каталог`,password:`Пароль доставки`,signIn:`Подтвердить доступ`,accessTitle:`Проверка доступа`,accessHelp:`Маршрут защищён. Подтвердите доступ, чтобы увидеть данные доставки.`,empty:`По этому пути нет доступных объектов.`,failed:`Маршрут недоступен или требуется авторизация. Данные доставки не были раскрыты.`,ready:`Маршрут готов`,manifest:`Манифест доставки`,confirmed:`Подтверждённая передача`,entryTypeFile:`Файл`,entryTypeDirectory:`Каталог`,itemSize:`Байт`}};function it(e,t){return rt[e][t]}Xe();var at=class extends K{constructor(...e){super(...e),this.locale=Ae(),this.failed=!1,this.csrf=``}static{this.properties={locale:{state:!0},metadata:{state:!0},failed:{state:!0},csrf:{state:!0}}}static{this.styles=[be,o`
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
  `]}connectedCallback(){super.connectedCallback(),this.theme=We(),this.refresh()}disconnectedCallback(){this.theme?.destroy(),super.disconnectedCallback()}async refresh(e=this.metadata?.path??``){this.failed=!1;try{this.metadata=await Qe(e)}catch{this.failed=!0,this.metadata=void 0}}async openDirectory(e,t){e.preventDefault(),await this.refresh(t)}async signIn(e){e.preventDefault();let t=e.currentTarget,n=new FormData(t).get(`password`)?.toString()??``;try{this.csrf=(await $e(n)).csrf,t.reset(),await this.refresh()}catch{this.failed=!0}}async sendFile(e){let t=e.currentTarget,n=t.files?.item(0);if(n)try{await et(n,this.csrf),t.value=``,await this.refresh()}catch{this.failed=!0}}setLocale(e){this.locale=e.detail}t(e){return it(this.locale,e)}entry(e){let t=tt(this.metadata.path,e.name);return e.type===`directory`?F`<li><courier-icon name="folder"></courier-icon><button class="link entry-name" @click=${e=>this.openDirectory(e,t)}>${e.name}</button><span class="size">${e.size} ${this.t(`itemSize`)}</span><a href=${$(t,!0)}>${this.t(`downloadArchive`)}</a></li>`:F`<li><courier-icon name="parcel"></courier-icon><span class="entry-name">${e.name}</span><span class="size">${e.size} ${this.t(`itemSize`)}</span><a href=${$(t)}>${this.t(`download`)}</a></li>`}render(){let e=this.metadata?.entries??[];return F`
      <main>
        <header>
          <courier-brand product=${this.t(`brandProduct`)}></courier-brand>
          <nav><courier-theme-selector .locale=${this.locale}></courier-theme-selector><courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector></nav>
        </header>
        <div class="workspace">
          <div class="operation-head"><div><span class="eyebrow">${this.t(`privateRoute`)}</span><h1>${this.t(`title`)}</h1></div>${this.metadata?F`<courier-status tone="signal">${this.t(`ready`)}</courier-status>`:L}</div>
          ${this.failed?F`
            <div class="access">
              <div class="access-copy">
                <span class="eyebrow">${this.t(`privateRoute`)}</span>
                <h2>${this.t(`accessTitle`)}</h2>
                <p class="muted">${this.t(`accessHelp`)}</p>
                <p class="error" role="alert">${this.t(`failed`)}</p>
                <form class="signin" @submit=${this.signIn}><label class="field"><span>${this.t(`password`)}</span><input name="password" type="password" autocomplete="current-password" placeholder=${this.t(`password`)}></label><courier-button type="submit" variant="primary">${this.t(`signIn`)}</courier-button></form>
                <courier-button @click=${this.refresh}>${this.t(`retry`)}</courier-button>
              </div>
              <div class="access-art"><courier-mascot alt="" .source=${Ze}></courier-mascot></div>
            </div>
          `:L}
          ${this.metadata?F`
            <div class="route-overview"><span class="label">${this.t(`confirmed`)}</span><courier-route source="sender" destination=${this.metadata.name}></courier-route></div>
            <section class="delivery-panel">
              <div class="delivery-title"><div><span class="eyebrow">${this.t(`manifest`)}</span><h2>${this.metadata.name}</h2></div><courier-status tone="signal">${this.t(`ready`)}</courier-status></div>
              ${this.metadata.type===`upload`?F`
                <div class="upload-zone"><h2>${this.t(`uploadTitle`)}</h2><p class="muted">${this.t(`uploadHelp`)}</p><label class="courier-file-action"><courier-icon name="upload"></courier-icon><span>${this.t(`upload`)}</span><input type="file" @change=${this.sendFile}></label></div>
              `:F`
                ${this.metadata.type===`file`?F`<div class="toolbar"><courier-icon name="download"></courier-icon><a href=${$(this.metadata.path)}>${this.t(`download`)}</a></div>`:F`
                  <div class="toolbar"><a href=${$(this.metadata.path,!0)}>${this.t(`downloadAll`)}</a>${this.metadata.path?F`<button class="link" @click=${e=>this.openDirectory(e,nt(this.metadata.path))}>${this.t(`up`)}</button>`:L}</div>
                  ${e.length===0?F`<p class="muted">${this.t(`empty`)}</p>`:F`<ul>${e.map(e=>this.entry(e))}</ul>`}
                `}
              `}
            </section>
          `:this.failed?L:F`<div class="loading"><courier-status>${this.t(`loading`)}</courier-status></div>`}
        </div>
      </main>
    `}};customElements.get(`courier-data-app`)||customElements.define(`courier-data-app`,at);