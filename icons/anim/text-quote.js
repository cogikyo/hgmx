import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bar = el.querySelector('path[d="M3 12v7"]');
  animate(bar, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid });

  const quoted = el.querySelectorAll('path[d^="M21 "]');
  animate(quoted, { translateX: [-3, 1, 0] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
