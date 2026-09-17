import { vectorMarkSource } from "@courier/ui";
import "./app";

const favicon = document.createElement("link");
favicon.rel = "icon";
favicon.type = "image/svg+xml";
favicon.href = vectorMarkSource;
document.head.append(favicon);
