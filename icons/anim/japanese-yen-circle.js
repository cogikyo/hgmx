import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = [el.querySelector('path[d="m9 7 3 3v8"]'), el.querySelector('path[d="m12 10 3-3"]')];
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d="M9 11h6"]'), el.querySelector('path[d="M9 15h6"]')];
  animate(bars, { scaleX: [0, 1.25, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
