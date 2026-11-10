import { animate, τmin } from "../motion.js";

export default function (el) {
  const f = el.querySelectorAll('path[d^="M8 20V8"], path[d="M6 12h4"], path[d="M6 20h4"]');
  animate(f, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3 });
  const i = el.querySelectorAll('path[d="M14 12h2v8"], path[d="M14 20h4"]');
  animate(i, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}
