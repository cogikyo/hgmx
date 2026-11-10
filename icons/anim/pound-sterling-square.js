import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="M10 17"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d="M8 13h5"]'), el.querySelector('path[d="M8 17h7"]')];
  animate(bars, { scaleX: [0, 1.25, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
