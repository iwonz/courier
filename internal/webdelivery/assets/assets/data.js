(function(){let e=document.createElement(`link`).relList;if(e&&e.supports&&e.supports(`modulepreload`))return;for(let e of document.querySelectorAll(`link[rel="modulepreload"]`))n(e);new MutationObserver(e=>{for(let t of e)if(t.type===`childList`)for(let e of t.addedNodes)e.tagName===`LINK`&&e.rel===`modulepreload`&&n(e)}).observe(document,{childList:!0,subtree:!0});function t(e){let t={};return e.integrity&&(t.integrity=e.integrity),e.referrerPolicy&&(t.referrerPolicy=e.referrerPolicy),t.credentials=e.crossOrigin===`use-credentials`?`include`:e.crossOrigin===`anonymous`?`omit`:`same-origin`,t}function n(e){if(e.ep)return;e.ep=!0;let n=t(e);fetch(e.href,n)}})();var e=globalThis,t=e.ShadowRoot&&(e.ShadyCSS===void 0||e.ShadyCSS.nativeShadow)&&`adoptedStyleSheets`in Document.prototype&&`replace`in CSSStyleSheet.prototype,n=Symbol(),r=new WeakMap,i=class{constructor(e,t,r){if(this._$cssResult$=!0,r!==n)throw Error("CSSResult is not constructable. Use `unsafeCSS` or `css` instead.");this.cssText=e,this.t=t}get styleSheet(){let e=this.o,n=this.t;if(t&&e===void 0){let t=n!==void 0&&n.length===1;t&&(e=r.get(n)),e===void 0&&((this.o=e=new CSSStyleSheet).replaceSync(this.cssText),t&&r.set(n,e))}return e}toString(){return this.cssText}},a=e=>new i(typeof e==`string`?e:e+``,void 0,n),o=(e,...t)=>new i(e.length===1?e[0]:t.reduce((t,n,r)=>t+(e=>{if(!0===e._$cssResult$)return e.cssText;if(typeof e==`number`)return e;throw Error(`Value passed to 'css' function must be a 'css' function result: `+e+`. Use 'unsafeCSS' to pass non-literal values, but take care to ensure page security.`)})(n)+e[r+1],e[0]),e,n),s=(n,r)=>{if(t)n.adoptedStyleSheets=r.map(e=>e instanceof CSSStyleSheet?e:e.styleSheet);else for(let t of r){let r=document.createElement(`style`),i=e.litNonce;i!==void 0&&r.setAttribute(`nonce`,i),r.textContent=t.cssText,n.appendChild(r)}},c=t?e=>e:e=>e instanceof CSSStyleSheet?(e=>{let t=``;for(let n of e.cssRules)t+=n.cssText;return a(t)})(e):e,{is:l,defineProperty:u,getOwnPropertyDescriptor:d,getOwnPropertyNames:ee,getOwnPropertySymbols:te,getPrototypeOf:ne}=Object,f=globalThis,p=f.trustedTypes,re=p?p.emptyScript:``,ie=f.reactiveElementPolyfillSupport,m=(e,t)=>e,h={toAttribute(e,t){switch(t){case Boolean:e=e?re:null;break;case Object:case Array:e=e==null?e:JSON.stringify(e)}return e},fromAttribute(e,t){let n=e;switch(t){case Boolean:n=e!==null;break;case Number:n=e===null?null:Number(e);break;case Object:case Array:try{n=JSON.parse(e)}catch{n=null}}return n}},g=(e,t)=>!l(e,t),_={attribute:!0,type:String,converter:h,reflect:!1,useDefault:!1,hasChanged:g};Symbol.metadata??=Symbol(`metadata`),f.litPropertyMetadata??=new WeakMap;var v=class extends HTMLElement{static addInitializer(e){this._$Ei(),(this.l??=[]).push(e)}static get observedAttributes(){return this.finalize(),this._$Eh&&[...this._$Eh.keys()]}static createProperty(e,t=_){if(t.state&&(t.attribute=!1),this._$Ei(),this.prototype.hasOwnProperty(e)&&((t=Object.create(t)).wrapped=!0),this.elementProperties.set(e,t),!t.noAccessor){let n=Symbol(),r=this.getPropertyDescriptor(e,n,t);r!==void 0&&u(this.prototype,e,r)}}static getPropertyDescriptor(e,t,n){let{get:r,set:i}=d(this.prototype,e)??{get(){return this[t]},set(e){this[t]=e}};return{get:r,set(t){let a=r?.call(this);i?.call(this,t),this.requestUpdate(e,a,n)},configurable:!0,enumerable:!0}}static getPropertyOptions(e){return this.elementProperties.get(e)??_}static _$Ei(){if(this.hasOwnProperty(m(`elementProperties`)))return;let e=ne(this);e.finalize(),e.l!==void 0&&(this.l=[...e.l]),this.elementProperties=new Map(e.elementProperties)}static finalize(){if(this.hasOwnProperty(m(`finalized`)))return;if(this.finalized=!0,this._$Ei(),this.hasOwnProperty(m(`properties`))){let e=this.properties,t=[...ee(e),...te(e)];for(let n of t)this.createProperty(n,e[n])}let e=this[Symbol.metadata];if(e!==null){let t=litPropertyMetadata.get(e);if(t!==void 0)for(let[e,n]of t)this.elementProperties.set(e,n)}this._$Eh=new Map;for(let[e,t]of this.elementProperties){let n=this._$Eu(e,t);n!==void 0&&this._$Eh.set(n,e)}this.elementStyles=this.finalizeStyles(this.styles)}static finalizeStyles(e){let t=[];if(Array.isArray(e)){let n=new Set(e.flat(1/0).reverse());for(let e of n)t.unshift(c(e))}else e!==void 0&&t.push(c(e));return t}static _$Eu(e,t){let n=t.attribute;return!1===n?void 0:typeof n==`string`?n:typeof e==`string`?e.toLowerCase():void 0}constructor(){super(),this._$Ep=void 0,this.isUpdatePending=!1,this.hasUpdated=!1,this._$Em=null,this._$Ev()}_$Ev(){this._$ES=new Promise(e=>this.enableUpdating=e),this._$AL=new Map,this._$E_(),this.requestUpdate(),this.constructor.l?.forEach(e=>e(this))}addController(e){(this._$EO??=new Set).add(e),this.renderRoot!==void 0&&this.isConnected&&e.hostConnected?.()}removeController(e){this._$EO?.delete(e)}_$E_(){let e=new Map,t=this.constructor.elementProperties;for(let n of t.keys())this.hasOwnProperty(n)&&(e.set(n,this[n]),delete this[n]);e.size>0&&(this._$Ep=e)}createRenderRoot(){let e=this.shadowRoot??this.attachShadow(this.constructor.shadowRootOptions);return s(e,this.constructor.elementStyles),e}connectedCallback(){this.renderRoot??=this.createRenderRoot(),this.enableUpdating(!0),this._$EO?.forEach(e=>e.hostConnected?.())}enableUpdating(e){}disconnectedCallback(){this._$EO?.forEach(e=>e.hostDisconnected?.())}attributeChangedCallback(e,t,n){this._$AK(e,n)}_$ET(e,t){let n=this.constructor.elementProperties.get(e),r=this.constructor._$Eu(e,n);if(r!==void 0&&!0===n.reflect){let i=(n.converter?.toAttribute===void 0?h:n.converter).toAttribute(t,n.type);this._$Em=e,i==null?this.removeAttribute(r):this.setAttribute(r,i),this._$Em=null}}_$AK(e,t){let n=this.constructor,r=n._$Eh.get(e);if(r!==void 0&&this._$Em!==r){let e=n.getPropertyOptions(r),i=typeof e.converter==`function`?{fromAttribute:e.converter}:e.converter?.fromAttribute===void 0?h:e.converter;this._$Em=r;let a=i.fromAttribute(t,e.type);this[r]=a??this._$Ej?.get(r)??a,this._$Em=null}}requestUpdate(e,t,n,r=!1,i){if(e!==void 0){let a=this.constructor;if(!1===r&&(i=this[e]),n??=a.getPropertyOptions(e),!((n.hasChanged??g)(i,t)||n.useDefault&&n.reflect&&i===this._$Ej?.get(e)&&!this.hasAttribute(a._$Eu(e,n))))return;this.C(e,t,n)}!1===this.isUpdatePending&&(this._$ES=this._$EP())}C(e,t,{useDefault:n,reflect:r,wrapped:i},a){n&&!(this._$Ej??=new Map).has(e)&&(this._$Ej.set(e,a??t??this[e]),!0!==i||a!==void 0)||(this._$AL.has(e)||(this.hasUpdated||n||(t=void 0),this._$AL.set(e,t)),!0===r&&this._$Em!==e&&(this._$Eq??=new Set).add(e))}async _$EP(){this.isUpdatePending=!0;try{await this._$ES}catch(e){Promise.reject(e)}let e=this.scheduleUpdate();return e!=null&&await e,!this.isUpdatePending}scheduleUpdate(){return this.performUpdate()}performUpdate(){if(!this.isUpdatePending)return;if(!this.hasUpdated){if(this.renderRoot??=this.createRenderRoot(),this._$Ep){for(let[e,t]of this._$Ep)this[e]=t;this._$Ep=void 0}let e=this.constructor.elementProperties;if(e.size>0)for(let[t,n]of e){let{wrapped:e}=n,r=this[t];!0!==e||this._$AL.has(t)||r===void 0||this.C(t,void 0,n,r)}}let e=!1,t=this._$AL;try{e=this.shouldUpdate(t),e?(this.willUpdate(t),this._$EO?.forEach(e=>e.hostUpdate?.()),this.update(t)):this._$EM()}catch(t){throw e=!1,this._$EM(),t}e&&this._$AE(t)}willUpdate(e){}_$AE(e){this._$EO?.forEach(e=>e.hostUpdated?.()),this.hasUpdated||(this.hasUpdated=!0,this.firstUpdated(e)),this.updated(e)}_$EM(){this._$AL=new Map,this.isUpdatePending=!1}get updateComplete(){return this.getUpdateComplete()}getUpdateComplete(){return this._$ES}shouldUpdate(e){return!0}update(e){this._$Eq&&=this._$Eq.forEach(e=>this._$ET(e,this[e])),this._$EM()}updated(e){}firstUpdated(e){}};v.elementStyles=[],v.shadowRootOptions={mode:`open`},v[m(`elementProperties`)]=new Map,v[m(`finalized`)]=new Map,ie?.({ReactiveElement:v}),(f.reactiveElementVersions??=[]).push(`2.1.2`);var y=globalThis,b=e=>e,x=y.trustedTypes,S=x?x.createPolicy(`lit-html`,{createHTML:e=>e}):void 0,ae=`$lit$`,C=`lit$${Math.random().toFixed(9).slice(2)}$`,oe=`?`+C,se=`<${oe}>`,w=document,T=()=>w.createComment(``),E=e=>e===null||typeof e!=`object`&&typeof e!=`function`,D=Array.isArray,ce=e=>D(e)||typeof e?.[Symbol.iterator]==`function`,O=`[ 	
\f\r]`,k=/<(?:(!--|\/[^a-zA-Z])|(\/?[a-zA-Z][^>\s]*)|(\/?$))/g,A=/-->/g,j=/>/g,M=RegExp(`>|${O}(?:([^\\s"'>=/]+)(${O}*=${O}*(?:[^ \t\n\f\r"'\`<>=]|("|')|))|$)`,`g`),N=/'/g,P=/"/g,F=/^(?:script|style|textarea|title)$/i,I=(e=>(t,...n)=>({_$litType$:e,strings:t,values:n}))(1),L=Symbol.for(`lit-noChange`),R=Symbol.for(`lit-nothing`),z=new WeakMap,B=w.createTreeWalker(w,129);function V(e,t){if(!D(e)||!e.hasOwnProperty(`raw`))throw Error(`invalid template strings array`);return S===void 0?t:S.createHTML(t)}var le=(e,t)=>{let n=e.length-1,r=[],i,a=t===2?`<svg>`:t===3?`<math>`:``,o=k;for(let t=0;t<n;t++){let n=e[t],s,c,l=-1,u=0;for(;u<n.length&&(o.lastIndex=u,c=o.exec(n),c!==null);)u=o.lastIndex,o===k?c[1]===`!--`?o=A:c[1]===void 0?c[2]===void 0?c[3]!==void 0&&(o=M):(F.test(c[2])&&(i=RegExp(`</`+c[2],`g`)),o=M):o=j:o===M?c[0]===`>`?(o=i??k,l=-1):c[1]===void 0?l=-2:(l=o.lastIndex-c[2].length,s=c[1],o=c[3]===void 0?M:c[3]===`"`?P:N):o===P||o===N?o=M:o===A||o===j?o=k:(o=M,i=void 0);let d=o===M&&e[t+1].startsWith(`/>`)?` `:``;a+=o===k?n+se:l>=0?(r.push(s),n.slice(0,l)+ae+n.slice(l)+C+d):n+C+(l===-2?t:d)}return[V(e,a+(e[n]||`<?>`)+(t===2?`</svg>`:t===3?`</math>`:``)),r]},H=class e{constructor({strings:t,_$litType$:n},r){let i;this.parts=[];let a=0,o=0,s=t.length-1,c=this.parts,[l,u]=le(t,n);if(this.el=e.createElement(l,r),B.currentNode=this.el.content,n===2||n===3){let e=this.el.content.firstChild;e.replaceWith(...e.childNodes)}for(;(i=B.nextNode())!==null&&c.length<s;){if(i.nodeType===1){if(i.hasAttributes())for(let e of i.getAttributeNames())if(e.endsWith(ae)){let t=u[o++],n=i.getAttribute(e).split(C),r=/([.?@])?(.*)/.exec(t);c.push({type:1,index:a,name:r[2],strings:n,ctor:r[1]===`.`?de:r[1]===`?`?fe:r[1]===`@`?pe:G}),i.removeAttribute(e)}else e.startsWith(C)&&(c.push({type:6,index:a}),i.removeAttribute(e));if(F.test(i.tagName)){let e=i.textContent.split(C),t=e.length-1;if(t>0){i.textContent=x?x.emptyScript:``;for(let n=0;n<t;n++)i.append(e[n],T()),B.nextNode(),c.push({type:2,index:++a});i.append(e[t],T())}}}else if(i.nodeType===8){if(i.data===oe)c.push({type:2,index:a});else{let e=-1;for(;(e=i.data.indexOf(C,e+1))!==-1;)c.push({type:7,index:a}),e+=C.length-1}}a++}}static createElement(e,t){let n=w.createElement(`template`);return n.innerHTML=e,n}};function U(e,t,n=e,r){if(t===L)return t;let i=r===void 0?n._$Cl:n._$Co?.[r],a=E(t)?void 0:t._$litDirective$;return i?.constructor!==a&&(i?._$AO?.(!1),a===void 0?i=void 0:(i=new a(e),i._$AT(e,n,r)),r===void 0?n._$Cl=i:(n._$Co??=[])[r]=i),i!==void 0&&(t=U(e,i._$AS(e,t.values),i,r)),t}var ue=class{constructor(e,t){this._$AV=[],this._$AN=void 0,this._$AD=e,this._$AM=t}get parentNode(){return this._$AM.parentNode}get _$AU(){return this._$AM._$AU}u(e){let{el:{content:t},parts:n}=this._$AD,r=(e?.creationScope??w).importNode(t,!0);B.currentNode=r;let i=B.nextNode(),a=0,o=0,s=n[0];for(;s!==void 0;){if(a===s.index){let t;s.type===2?t=new W(i,i.nextSibling,this,e):s.type===1?t=new s.ctor(i,s.name,s.strings,this,e):s.type===6&&(t=new me(i,this,e)),this._$AV.push(t),s=n[++o]}a!==s?.index&&(i=B.nextNode(),a++)}return B.currentNode=w,r}p(e){let t=0;for(let n of this._$AV)n!==void 0&&(n.strings===void 0?n._$AI(e[t]):(n._$AI(e,n,t),t+=n.strings.length-2)),t++}},W=class e{get _$AU(){return this._$AM?._$AU??this._$Cv}constructor(e,t,n,r){this.type=2,this._$AH=R,this._$AN=void 0,this._$AA=e,this._$AB=t,this._$AM=n,this.options=r,this._$Cv=r?.isConnected??!0}get parentNode(){let e=this._$AA.parentNode,t=this._$AM;return t!==void 0&&e?.nodeType===11&&(e=t.parentNode),e}get startNode(){return this._$AA}get endNode(){return this._$AB}_$AI(e,t=this){e=U(this,e,t),E(e)?e===R||e==null||e===``?(this._$AH!==R&&this._$AR(),this._$AH=R):e!==this._$AH&&e!==L&&this._(e):e._$litType$===void 0?e.nodeType===void 0?ce(e)?this.k(e):this._(e):this.T(e):this.$(e)}O(e){return this._$AA.parentNode.insertBefore(e,this._$AB)}T(e){this._$AH!==e&&(this._$AR(),this._$AH=this.O(e))}_(e){this._$AH!==R&&E(this._$AH)?this._$AA.nextSibling.data=e:this.T(w.createTextNode(e)),this._$AH=e}$(e){let{values:t,_$litType$:n}=e,r=typeof n==`number`?this._$AC(e):(n.el===void 0&&(n.el=H.createElement(V(n.h,n.h[0]),this.options)),n);if(this._$AH?._$AD===r)this._$AH.p(t);else{let e=new ue(r,this),n=e.u(this.options);e.p(t),this.T(n),this._$AH=e}}_$AC(e){let t=z.get(e.strings);return t===void 0&&z.set(e.strings,t=new H(e)),t}k(t){D(this._$AH)||(this._$AH=[],this._$AR());let n=this._$AH,r,i=0;for(let a of t)i===n.length?n.push(r=new e(this.O(T()),this.O(T()),this,this.options)):r=n[i],r._$AI(a),i++;i<n.length&&(this._$AR(r&&r._$AB.nextSibling,i),n.length=i)}_$AR(e=this._$AA.nextSibling,t){for(this._$AP?.(!1,!0,t);e!==this._$AB;){let t=b(e).nextSibling;b(e).remove(),e=t}}setConnected(e){this._$AM===void 0&&(this._$Cv=e,this._$AP?.(e))}},G=class{get tagName(){return this.element.tagName}get _$AU(){return this._$AM._$AU}constructor(e,t,n,r,i){this.type=1,this._$AH=R,this._$AN=void 0,this.element=e,this.name=t,this._$AM=r,this.options=i,n.length>2||n[0]!==``||n[1]!==``?(this._$AH=Array(n.length-1).fill(new String),this.strings=n):this._$AH=R}_$AI(e,t=this,n,r){let i=this.strings,a=!1;if(i===void 0)e=U(this,e,t,0),a=!E(e)||e!==this._$AH&&e!==L,a&&(this._$AH=e);else{let r=e,o,s;for(e=i[0],o=0;o<i.length-1;o++)s=U(this,r[n+o],t,o),s===L&&(s=this._$AH[o]),a||=!E(s)||s!==this._$AH[o],s===R?e=R:e!==R&&(e+=(s??``)+i[o+1]),this._$AH[o]=s}a&&!r&&this.j(e)}j(e){e===R?this.element.removeAttribute(this.name):this.element.setAttribute(this.name,e??``)}},de=class extends G{constructor(){super(...arguments),this.type=3}j(e){this.element[this.name]=e===R?void 0:e}},fe=class extends G{constructor(){super(...arguments),this.type=4}j(e){this.element.toggleAttribute(this.name,!!e&&e!==R)}},pe=class extends G{constructor(e,t,n,r,i){super(e,t,n,r,i),this.type=5}_$AI(e,t=this){if((e=U(this,e,t,0)??R)===L)return;let n=this._$AH,r=e===R&&n!==R||e.capture!==n.capture||e.once!==n.once||e.passive!==n.passive,i=e!==R&&(n===R||r);r&&this.element.removeEventListener(this.name,this,n),i&&this.element.addEventListener(this.name,this,e),this._$AH=e}handleEvent(e){typeof this._$AH==`function`?this._$AH.call(this.options?.host??this.element,e):this._$AH.handleEvent(e)}},me=class{constructor(e,t,n){this.element=e,this.type=6,this._$AN=void 0,this._$AM=t,this.options=n}get _$AU(){return this._$AM._$AU}_$AI(e){U(this,e)}},he=y.litHtmlPolyfillSupport;he?.(H,W),(y.litHtmlVersions??=[]).push(`3.3.3`);var ge=(e,t,n)=>{let r=n?.renderBefore??t,i=r._$litPart$;if(i===void 0){let e=n?.renderBefore??null;r._$litPart$=i=new W(t.insertBefore(T(),e),e,void 0,n??{})}return i._$AI(e),i},K=globalThis,q=class extends v{constructor(){super(...arguments),this.renderOptions={host:this},this._$Do=void 0}createRenderRoot(){let e=super.createRenderRoot();return this.renderOptions.renderBefore??=e.firstChild,e}update(e){let t=this.render();this.hasUpdated||(this.renderOptions.isConnected=this.isConnected),super.update(e),this._$Do=ge(t,this.renderRoot,this.renderOptions)}connectedCallback(){super.connectedCallback(),this._$Do?.setConnected(!0)}disconnectedCallback(){super.disconnectedCallback(),this._$Do?.setConnected(!1)}render(){return L}};q._$litElement$=!0,q.finalized=!0,K.litElementHydrateSupport?.({LitElement:q});var _e=K.litElementPolyfillSupport;_e?.({LitElement:q}),(K.litElementVersions??=[]).push(`4.2.2`);var ve=o`
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
`;var ye=o`
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
`,be=class extends q{constructor(...e){super(...e),this.disabled=!1,this.type=`button`,this.variant=`secondary`}static{this.properties={disabled:{type:Boolean,reflect:!0},type:{type:String,reflect:!0},variant:{type:String,reflect:!0}}}static{this.styles=[ve,o`
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
  `]}render(){return I`<button type=${this.type} ?disabled=${this.disabled}><slot></slot></button>`}},xe=`/assets/data.png`,Se=class extends q{constructor(...e){super(...e),this.product=``}static{this.properties={product:{type:String}}}static{this.styles=o`
    :host { display: inline-flex; min-width: 0; color: var(--courier-color-text, #151714); font-family: var(--courier-font-sans, sans-serif); }
    .lockup { display: inline-flex; min-width: 0; align-items: center; gap: 0.625rem; color: inherit; }
    img { flex: 0 0 auto; width: 2.5rem; height: 2.5rem; object-fit: contain; filter: drop-shadow(0 0.3rem 0.45rem rgb(16 18 15 / 0.16)); }
    .words { display: grid; line-height: 1; }
    strong { font-family: var(--courier-font-display, sans-serif); font-size: 1.125rem; font-weight: 850; letter-spacing: -0.04em; }
    small { margin-top: 0.25rem; color: var(--courier-color-muted, #596054); font-family: var(--courier-font-mono, monospace); font-size: 0.625rem; font-weight: 700; letter-spacing: 0.09em; text-transform: uppercase; }
  `}render(){return I`<span class="lockup">
      <img part="mark" src=${xe} alt="" decoding="async">
      <span class="words" part="words"><strong>Courier</strong><small part="product">${this.product}</small></span>
    </span>`}},Ce=class extends q{constructor(...e){super(...e),this.alt=``,this.eager=!1,this.mobileSource=``,this.source=``}static{this.properties={alt:{type:String},eager:{type:Boolean},mobileSource:{type:String,attribute:`mobile-source`},source:{type:String}}}static{this.styles=o`
    :host { display: block; }
    picture { display: contents; }
    img { display: block; width: 100%; height: auto; filter: drop-shadow(0 1.5rem 2rem rgb(16 18 15 / 0.18)); }
  `}render(){return I`<picture>${this.mobileSource?I`<source media="(max-width: 44rem)" srcset=${this.mobileSource}>`:``}<img part="image" src=${this.source} alt=${this.alt} decoding="async" loading=${this.eager?`eager`:`lazy`} fetchpriority=${this.eager?`high`:`auto`}></picture>`}},we=class extends q{constructor(...e){super(...e),this.tone=`neutral`}static{this.properties={tone:{type:String,reflect:!0}}}static{this.styles=o`
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
  `}render(){return I`<i aria-hidden="true"></i><slot></slot>`}},Te=class extends q{constructor(...e){super(...e),this.source=``,this.destination=``}static{this.properties={source:{type:String},destination:{type:String}}}static{this.styles=o`
    :host { display: grid; color: var(--courier-color-text, #151714); font-family: var(--courier-font-mono, monospace); }
    .route { display: grid; grid-template-columns: minmax(0, 1fr) minmax(3rem, 0.55fr) minmax(0, 1fr); align-items: center; gap: 0.65rem; }
    .node { overflow: hidden; padding: 0.65rem 0.75rem; border: 1px solid var(--courier-color-border, #c8cdbf); border-radius: var(--courier-radius-sm, 0.25rem); background: var(--courier-color-surface, #fafbf3); font-size: 0.75rem; text-overflow: ellipsis; white-space: nowrap; }
    .line { position: relative; height: 1px; color: var(--courier-color-border-strong, #8e9587); background: currentColor; }
    .line::before { content: ""; position: absolute; top: -0.2rem; left: 0; width: 0.45rem; height: 0.45rem; border-radius: 50%; background: var(--courier-color-accent, #d4ff45); }
    .line::after { content: ""; position: absolute; top: -0.22rem; right: 0; width: 0.4rem; height: 0.4rem; border-top: 1px solid currentColor; border-right: 1px solid currentColor; transform: rotate(45deg); }
  `}render(){return I`<div class="route"><span class="node">${this.source}</span><span class="line" aria-hidden="true"></span><span class="node">${this.destination}</span></div>`}},Ee={en:{"theme.label":`Theme`,"theme.system":`System`,"theme.light":`Light`,"theme.dark":`Dark`,"locale.label":`Language`,"locale.en":`English`,"locale.ru":`Russian`,"progress.label":`Delivery progress`,"action.cancel":`Cancel`,"action.close":`Close`},ru:{"theme.label":`Тема`,"theme.system":`Системная`,"theme.light":`Светлая`,"theme.dark":`Тёмная`,"locale.label":`Язык`,"locale.en":`Английский`,"locale.ru":`Русский`,"progress.label":`Ход доставки`,"action.cancel":`Отмена`,"action.close":`Закрыть`}},De=Object.keys(Ee),J=`courier.locale`;function Y(e){if(!e)return;let t=e.toLowerCase().split(`-`)[0];return De.includes(t)?t:void 0}function Oe(e){for(let t of e){let e=Y(t);if(e)return e}return`en`}function ke(e,t){if(e)try{let t=Y(e.getItem(J));if(t)return t}catch{}return Oe(t)}function Ae(e,t){if(e)try{e.setItem(J,t)}catch{}}function X(e,t){return Ee[Y(e)??`en`][t]}function je(){let e;try{e=globalThis.localStorage}catch{e=void 0}let t=globalThis.navigator?.languages??[];return ke(e,t)}var Me=class extends q{constructor(...e){super(...e),this.locale=`en`}static{this.properties={locale:{type:String}}}connectedCallback(){super.connectedCallback(),this.locale=je()}change(e){let t=Y(e.detail)??`en`;this.locale=t;let n;try{n=globalThis.localStorage}catch{n=void 0}Ae(n,t),this.dispatchEvent(new CustomEvent(`courier-locale-change`,{detail:t,bubbles:!0,composed:!0}))}render(){return I`<courier-segmented-control
      icon-only
      .label=${X(this.locale,`locale.label`)}
      .value=${this.locale}
      .options=${[{value:`en`,label:X(this.locale,`locale.en`),symbol:`🇬🇧`},{value:`ru`,label:X(this.locale,`locale.ru`),symbol:`🇷🇺`}]}
      data-storage-key=${J}
      @courier-segment-change=${this.change}
    ></courier-segmented-control>`}},Ne=class extends q{constructor(...e){super(...e),this.heading=``}static{this.properties={heading:{type:String}}}static{this.styles=o`
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
  `}render(){return this.heading?I`<section aria-labelledby="courier-panel-heading">
      <h2 id="courier-panel-heading">${this.heading}</h2>
      <slot></slot>
    </section>`:I`<section><slot></slot></section>`}};function Pe(e,t){return!Number.isFinite(e)||!Number.isFinite(t)||t<=0?0:Math.min(1,Math.max(0,e/t))}var Fe=class extends q{constructor(...e){super(...e),this.value=0,this.total=0,this.label=``,this.locale=`en`}static{this.properties={value:{type:Number},total:{type:Number},label:{type:String},locale:{type:String}}}static{this.styles=o`
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
  `}render(){let e=Pe(this.value,this.total),t=this.label||X(this.locale,`progress.label`),n=Number.isFinite(this.total)&&this.total>0?this.total:0;return I`<div
      class="track"
      role="progressbar"
      aria-label=${t}
      aria-valuemin="0"
      aria-valuemax=${n}
      aria-valuenow=${Number.isFinite(this.value)?Math.max(0,Math.min(this.value,n)):0}
    ><div class="fill" style=${`transform: scaleX(${e})`}></div></div>
    <output>${Math.round(e*100)}%</output>`}};function Ie(e,t,n){if(!(n<=0))switch(e){case`ArrowLeft`:case`ArrowUp`:return(t-1+n)%n;case`ArrowRight`:case`ArrowDown`:return(t+1)%n;case`Home`:return 0;case`End`:return n-1;default:return}}var Le=class extends q{constructor(...e){super(...e),this.label=``,this.options=[],this.value=``,this.iconOnly=!1}static{this.properties={label:{type:String},options:{attribute:!1},value:{type:String},iconOnly:{type:Boolean,attribute:`icon-only`,reflect:!0}}}static{this.styles=o`
    :host { display: block; min-width: 0; color: var(--courier-color-text, #151714); font-family: var(--courier-font-sans, sans-serif); }
    fieldset { min-width: 0; margin: 0; padding: 0; border: 0; }
    legend { margin: 0 0 0.25rem; padding: 0; color: var(--courier-color-muted, #596054); font-family: var(--courier-font-mono, monospace); font-size: 0.625rem; font-weight: 750; letter-spacing: 0.08em; line-height: 1; text-transform: uppercase; }
    legend.sr-only { position: absolute; width: 1px; height: 1px; margin: -1px; padding: 0; overflow: hidden; clip: rect(0 0 0 0); clip-path: inset(50%); white-space: nowrap; }
    .segments { display: inline-grid; max-width: 100%; grid-auto-columns: minmax(0, auto); grid-auto-flow: column; gap: 2px; padding: 2px; border: 1px solid var(--courier-color-border, #c8cdbf); border-radius: var(--courier-radius-md, 0.625rem); background: color-mix(in srgb, var(--courier-color-field, #e7e9dc) 68%, transparent); box-shadow: inset 0 1px 2px rgb(16 18 15 / 0.07); }
    button { appearance: none; display: inline-flex; min-width: 0; min-height: 2.25rem; align-items: center; justify-content: center; gap: 0.35rem; padding: 0.45rem 0.68rem; overflow: hidden; border: 1px solid transparent; border-radius: calc(var(--courier-radius-md, 0.625rem) - 3px); color: var(--courier-color-muted, #596054); background: transparent; font: inherit; font-size: 0.75rem; font-weight: 780; line-height: 1; text-overflow: ellipsis; white-space: nowrap; cursor: pointer; transition: color var(--courier-duration, 160ms) var(--courier-ease, ease), background var(--courier-duration, 160ms) var(--courier-ease, ease), transform var(--courier-duration, 160ms) var(--courier-ease, ease); }
    :host([icon-only]) button { width: 2.25rem; padding: 0.45rem; }
    courier-icon { width: 1.05rem; height: 1.05rem; }
    .symbol { font-family: system-ui, sans-serif; font-size: 1.05rem; line-height: 1; }
    button:hover { color: var(--courier-color-text, #151714); background: color-mix(in srgb, var(--courier-color-surface-raised, #fff) 72%, transparent); }
    button.selected { border-color: color-mix(in srgb, var(--courier-color-accent, #d4ff45) 64%, var(--courier-color-border, #c8cdbf)); color: var(--courier-color-accent-ink, #151714); background: var(--courier-color-accent, #d4ff45); box-shadow: 0 1px 0 rgb(16 18 15 / 0.12); }
    button:active { transform: translateY(1px); }
    button:focus-visible { position: relative; z-index: 1; outline: 3px solid var(--courier-beak, #ff8758); outline-offset: 2px; }
  `}activate(e){let t=e.currentTarget.dataset.value??``;t&&t!==this.value&&(this.value=t,this.dispatchEvent(new CustomEvent(`courier-segment-change`,{detail:t,bubbles:!0,composed:!0})))}move(e){let t=[...this.renderRoot.querySelectorAll(`button`)],n=t.indexOf(e.currentTarget),r=Ie(e.key,n,t.length);if(r===void 0)return;e.preventDefault();let i=t[r];i.focus(),i.click()}render(){return I`<fieldset>
      <legend class=${this.iconOnly?`sr-only`:``}>${this.label}</legend>
      <div class="segments" role="radiogroup" aria-label=${this.label}>
        ${this.options.map(e=>{let t=e.value===this.value;return I`<button
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
          >${e.icon?I`<courier-icon name=${e.icon}></courier-icon>`:``}${e.symbol?I`<span class="symbol" aria-hidden="true">${e.symbol}</span>`:``}${this.iconOnly?``:I`<span>${e.label}</span>`}</button>`})}
      </div>
    </fieldset>`}},Re=[`system`,`light`,`dark`],ze=`courier.theme`;function Be(e){return Re.includes(e)?e:`system`}function Ve(e,t){return e===`system`?t?.matches?`dark`:`light`:e}function He(e){if(!e)return`system`;try{return Be(e.getItem(ze))}catch{return`system`}}function Ue(e,t){if(e)try{e.setItem(ze,t)}catch{}}var We=class{constructor(e,t,n,r){this.root=e,this.storage=t,this.media=n,this.onSystemChange=()=>this.apply(),this.preference=r??He(t),this.media?.addEventListener(`change`,this.onSystemChange),this.apply()}set(e){this.preference=e,Ue(this.storage,e),this.apply()}destroy(){this.media?.removeEventListener(`change`,this.onSystemChange)}apply(){this.root.dataset.courierTheme=Ve(this.preference,this.media),this.root.dataset.courierThemePreference=this.preference}};function Ge(){let e;try{e=globalThis.localStorage}catch{e=void 0}let t=globalThis.matchMedia?.(`(prefers-color-scheme: dark)`);return new We(document.documentElement,e,t)}var Ke=class extends q{constructor(...e){super(...e),this.preference=`system`,this.locale=`en`}static{this.properties={preference:{type:String},locale:{type:String}}}connectedCallback(){super.connectedCallback(),this.state=Ge(),this.preference=this.state.preference}disconnectedCallback(){this.state?.destroy(),super.disconnectedCallback()}change(e){let t=Be(e.detail);this.preference=t,this.state?.set(t),this.dispatchEvent(new CustomEvent(`courier-theme-change`,{detail:t,bubbles:!0,composed:!0}))}render(){return I`<courier-segmented-control
      icon-only
      .label=${X(this.locale,`theme.label`)}
      .value=${this.preference}
      .options=${[{value:`system`,label:X(this.locale,`theme.system`),icon:`system`},{value:`light`,label:X(this.locale,`theme.light`),icon:`sun`},{value:`dark`,label:X(this.locale,`theme.dark`),icon:`moon`}]}
      @courier-segment-change=${this.change}
    ></courier-segmented-control>`}},qe=[`apple`,`archive`,`copy`,`download`,`folder`,`github`,`homebrew`,`linux`,`moon`,`npm`,`package`,`parcel`,`pnpm`,`receipt`,`retry`,`route`,`scoop`,`server`,`shield`,`sun`,`system`,`terminal`,`upload`,`windows`,`yarn`],Je={apple:`M15 5c1-1 1-3 1-3-2 0-3 1-4 3m6 7c-1-2-2-3-4-3-1 0-2 1-3 1s-2-1-3-1c-3 0-5 3-5 6 0 4 3 8 5 8 1 0 2-1 3-1s2 1 3 1c2 0 4-3 5-6-2-1-3-2-3-4 0-2 1-3 2-4z`,archive:`M3 3h18v5H3zM5 8v13h14V8M9 12h6`,copy:`M8 3h13v13M3 8h13v13H3z`,download:`M12 3v13m-5-5 5 5 5-5M4 15v6h16v-6`,folder:`M3 6h7l2 3h9v12H3zM3 6V3h7l2 3h9v3`,github:`M9 19c-5 1-5-2-7-3m14 6v-3.6c0-1 .1-1.7-.4-2.2 3.2-.4 6.4-1.6 6.4-7.1 0-1.6-.6-3-1.7-4 .2-.5.7-2.3-.2-4.6 0 0-1.4-.5-4.7 1.7a16 16 0 0 0-8.6 0C6.4 1 5 1.5 5 1.5 4.1 3.8 4.6 5.6 4.8 6.1a7 7 0 0 0-1.7 4c0 5.5 3.2 6.7 6.4 7.1-.4.4-.8 1.1-.8 2.2V23`,homebrew:`M6 4h11l-1 15H8zM17 7h2a2 2 0 0 1 2 2v3a2 2 0 0 1-2 2h-2M5 22h13`,linux:`M12 2c-3 0-4 3-4 6-2 2-3 5-3 8l3-1 1 5 3-2 3 2 1-5 3 1c0-3-1-6-3-8 0-3-1-6-4-6zM9 8h.01M15 8h.01M10 11h4`,moon:`M20 16a8 8 0 0 1-12-10 8 8 0 1 0 12 10z`,npm:`M2 6h20v12H2zM6 15V9h5v6m0-6h4v6m0-6h3v6`,package:`m3 7 9-5 9 5v11l-9 4-9-4zM3 7l9 5 9-5M12 12v10M8 4l9 5`,parcel:`m3 7 9-5 9 5v11l-9 4-9-4zM3 7l9 5 9-5M12 12v10M8 4l9 5v5`,pnpm:`M3 3h5v5H3zM10 3h5v5h-5zM17 3h4v5h-4zM3 10h5v5H3zm7 0h5v5h-5zm7 0h4v5h-4zM10 17h5v4h-5zm7 0h4v4h-4z`,receipt:`M5 2h14v20l-3-2-4 2-4-2-3 2zM8 7h8M8 11h8m-8 5 2 2 5-4`,retry:`M3 10a9 9 0 1 1 1 7M3 3v7h7M12 7v5l3 2`,route:`M2 3h6v6H2zM16 15h6v6h-6zM11 6h8v6m-3-3 3 3 3-3M13 18H5v-6m-3 3 3-3 3 3`,scoop:`M5 8h14l-2 13H7zM4 8h16M8 8V5a4 4 0 0 1 8 0v3`,server:`M3 2h18v8H3zM3 14h18v8H3zM7 6h1m3 0h6M7 18h1m3 0h6M6 10v4m12-4v4`,shield:`m12 2 8 3v7c0 5-8 10-8 10S4 17 4 12V5zM8 11l3 3 5-6`,sun:`M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8zM12 2v3m0 14v3M4.9 4.9 7 7m10 10 2.1 2.1M2 12h3m14 0h3M4.9 19.1 7 17M17 7l2.1-2.1`,system:`M3 4h18v13H3zM8 21h8M12 17v4`,terminal:`M4 6l5 6-5 6m7 0h9`,upload:`M12 16V3m-5 5 5-5 5 5M4 15v6h16v-6`,windows:`M3 4l8-1v8H3zm10-1 8-1v9h-8zM3 13h8v8l-8-1zm10 0h8v9l-8-1z`,yarn:`M12 3a9 9 0 1 0 9 9M8 16c4-1 7-4 9-8m-8 1c3 1 5 4 5 8m-5-5c-1-3 0-5 2-6`};function Ye(e){return qe.includes(e)?e:`parcel`}var Xe={"courier-brand":Se,"courier-button":be,"courier-icon":class extends q{constructor(...e){super(...e),this.name=`parcel`,this.label=``}static{this.properties={name:{type:String},label:{type:String}}}static{this.styles=o`
    :host {
      display: inline-flex;
      width: 1.5rem;
      height: 1.5rem;
      color: currentColor;
    }
    svg { width: 100%; height: 100%; }
  `}render(){let e=Ye(this.name);return I`<svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="1.7"
      stroke-linecap="square"
      stroke-linejoin="miter"
      role=${this.label?`img`:`presentation`}
      aria-hidden=${this.label?`false`:`true`}
      aria-label=${this.label||void 0}
    ><path d=${Je[e]}></path></svg>`}},"courier-locale-selector":Me,"courier-mascot":Ce,"courier-panel":Ne,"courier-progress":Fe,"courier-route":Te,"courier-segmented-control":Le,"courier-status":we,"courier-theme-selector":Ke};function Ze(e=customElements){for(let[t,n]of Object.entries(Xe))e.get(t)||e.define(t,n)}var Qe=`/assets/data.webp`;function Z(e,t=``,n=globalThis.location.pathname){let r=n.endsWith(`/`)?n:`${n}/`,i=new URL(`api/v1/${e}`,globalThis.location.origin);return i.pathname=`${r}api/v1/${e}`,t&&i.searchParams.set(`path`,t),`${i.pathname}${i.search}`}async function Q(e){if(!e.ok)throw Error(`Courier request failed (${e.status})`);return e.json()}async function $e(e=``,t=globalThis.fetch){return Q(await t(Z(`meta`,e),{credentials:`same-origin`}))}async function et(e,t=globalThis.fetch){return Q(await t(Z(`session`),{method:`POST`,credentials:`same-origin`,headers:{"Content-Type":`application/json`},body:JSON.stringify({password:e})}))}async function tt(e,t,n=globalThis.fetch){let r=new FormData;r.append(`file`,e),await Q(await n(Z(`upload`),{method:`POST`,credentials:`same-origin`,headers:{"X-Courier-CSRF":t,"X-Courier-File-Size":String(e.size)},body:r}))}function $(e,t=!1){let n=new URL(Z(`download`,e),globalThis.location.origin);return t&&n.searchParams.set(`archive`,`tar.gz`),`${n.pathname}${n.search}`}function nt(e,t){return e?`${e}/${t}`:t}function rt(e){let t=e.lastIndexOf(`/`);return t<0?``:e.slice(0,t)}var it={en:{brandProduct:`Delivery terminal`,title:`Courier delivery`,privateRoute:`Private route`,loading:`Preparing the delivery route…`,retry:`Retry connection`,download:`Download file`,downloadArchive:`Download as archive`,downloadAll:`Download directory`,upload:`Choose a file`,uploadTitle:`Dispatch a file`,uploadHelp:`Select one file. Courier checks policy and reserves the final path before committing it.`,up:`Parent directory`,password:`Delivery password`,signIn:`Verify access`,accessTitle:`Identity check`,accessHelp:`This route is protected. Verify access to reveal its delivery metadata.`,empty:`No entries are available at this path.`,failed:`The route is unavailable or authorization is required. No delivery metadata was revealed.`,ready:`Route ready`,manifest:`Delivery manifest`,confirmed:`Verified handoff`,entryTypeFile:`File`,entryTypeDirectory:`Directory`,itemSize:`Bytes`},ru:{brandProduct:`Терминал доставки`,title:`Доставка Courier`,privateRoute:`Приватный маршрут`,loading:`Подготовка маршрута доставки…`,retry:`Повторить подключение`,download:`Скачать файл`,downloadArchive:`Скачать архивом`,downloadAll:`Скачать каталог`,upload:`Выбрать файл`,uploadTitle:`Отправить файл`,uploadHelp:`Выберите один файл. Courier проверит правила и зарезервирует конечный путь до фиксации.`,up:`Родительский каталог`,password:`Пароль доставки`,signIn:`Подтвердить доступ`,accessTitle:`Проверка доступа`,accessHelp:`Маршрут защищён. Подтвердите доступ, чтобы увидеть данные доставки.`,empty:`По этому пути нет доступных объектов.`,failed:`Маршрут недоступен или требуется авторизация. Данные доставки не были раскрыты.`,ready:`Маршрут готов`,manifest:`Манифест доставки`,confirmed:`Подтверждённая передача`,entryTypeFile:`Файл`,entryTypeDirectory:`Каталог`,itemSize:`Байт`}};function at(e,t){return it[e][t]}Ze();var ot=class extends q{constructor(...e){super(...e),this.locale=je(),this.failed=!1,this.csrf=``}static{this.properties={locale:{state:!0},metadata:{state:!0},failed:{state:!0},csrf:{state:!0}}}static{this.styles=[ye,o`
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
  `]}connectedCallback(){super.connectedCallback(),this.theme=Ge(),this.refresh()}disconnectedCallback(){this.theme?.destroy(),super.disconnectedCallback()}async refresh(e=this.metadata?.path??``){this.failed=!1;try{this.metadata=await $e(e)}catch{this.failed=!0,this.metadata=void 0}}async openDirectory(e,t){e.preventDefault(),await this.refresh(t)}async signIn(e){e.preventDefault();let t=e.currentTarget,n=new FormData(t).get(`password`)?.toString()??``;try{this.csrf=(await et(n)).csrf,t.reset(),await this.refresh()}catch{this.failed=!0}}async sendFile(e){let t=e.currentTarget,n=t.files?.item(0);if(n)try{await tt(n,this.csrf),t.value=``,await this.refresh()}catch{this.failed=!0}}setLocale(e){this.locale=e.detail}t(e){return at(this.locale,e)}entry(e){let t=nt(this.metadata.path,e.name);return e.type===`directory`?I`<li><courier-icon name="folder"></courier-icon><button class="link entry-name" @click=${e=>this.openDirectory(e,t)}>${e.name}</button><span class="size">${e.size} ${this.t(`itemSize`)}</span><a href=${$(t,!0)}>${this.t(`downloadArchive`)}</a></li>`:I`<li><courier-icon name="parcel"></courier-icon><span class="entry-name">${e.name}</span><span class="size">${e.size} ${this.t(`itemSize`)}</span><a href=${$(t)}>${this.t(`download`)}</a></li>`}render(){let e=this.metadata?.entries??[];return I`
      <main>
        <header>
          <courier-brand product=${this.t(`brandProduct`)}></courier-brand>
          <nav><courier-theme-selector .locale=${this.locale}></courier-theme-selector><courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector></nav>
        </header>
        <div class="workspace">
          <div class="operation-head"><div><span class="eyebrow">${this.t(`privateRoute`)}</span><h1>${this.t(`title`)}</h1></div>${this.metadata?I`<courier-status tone="signal">${this.t(`ready`)}</courier-status>`:R}</div>
          ${this.failed?I`
            <div class="access">
              <div class="access-copy">
                <span class="eyebrow">${this.t(`privateRoute`)}</span>
                <h2>${this.t(`accessTitle`)}</h2>
                <p class="muted">${this.t(`accessHelp`)}</p>
                <p class="error" role="alert">${this.t(`failed`)}</p>
                <form class="signin" @submit=${this.signIn}><label class="field"><span>${this.t(`password`)}</span><input name="password" type="password" autocomplete="current-password" placeholder=${this.t(`password`)}></label><courier-button type="submit" variant="primary">${this.t(`signIn`)}</courier-button></form>
                <courier-button @click=${this.refresh}>${this.t(`retry`)}</courier-button>
              </div>
              <div class="access-art"><courier-mascot alt="" .source=${Qe}></courier-mascot></div>
            </div>
          `:R}
          ${this.metadata?I`
            <div class="route-overview"><span class="label">${this.t(`confirmed`)}</span><courier-route source="sender" destination=${this.metadata.name}></courier-route></div>
            <section class="delivery-panel">
              <div class="delivery-title"><div><span class="eyebrow">${this.t(`manifest`)}</span><h2>${this.metadata.name}</h2></div><courier-status tone="signal">${this.t(`ready`)}</courier-status></div>
              ${this.metadata.type===`upload`?I`
                <div class="upload-zone"><h2>${this.t(`uploadTitle`)}</h2><p class="muted">${this.t(`uploadHelp`)}</p><label class="courier-file-action"><courier-icon name="upload"></courier-icon><span>${this.t(`upload`)}</span><input type="file" @change=${this.sendFile}></label></div>
              `:I`
                ${this.metadata.type===`file`?I`<div class="toolbar"><courier-icon name="download"></courier-icon><a href=${$(this.metadata.path)}>${this.t(`download`)}</a></div>`:I`
                  <div class="toolbar"><a href=${$(this.metadata.path,!0)}>${this.t(`downloadAll`)}</a>${this.metadata.path?I`<button class="link" @click=${e=>this.openDirectory(e,rt(this.metadata.path))}>${this.t(`up`)}</button>`:R}</div>
                  ${e.length===0?I`<p class="muted">${this.t(`empty`)}</p>`:I`<ul>${e.map(e=>this.entry(e))}</ul>`}
                `}
              `}
            </section>
          `:this.failed?R:I`<div class="loading"><courier-status>${this.t(`loading`)}</courier-status></div>`}
        </div>
      </main>
    `}};customElements.get(`courier-data-app`)||customElements.define(`courier-data-app`,ot);