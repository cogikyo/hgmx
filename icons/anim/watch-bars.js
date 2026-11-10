import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const face = el.querySelector("rect");
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const bars = [el.querySelector('path[d="M10 12v2"]'), el.querySelector('path[d="M14 10v4"]')];
  animate(bars, { scaleY: [1, 1.5, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
