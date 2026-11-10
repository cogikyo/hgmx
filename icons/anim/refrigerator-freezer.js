import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const fridge = el.querySelector("rect");
  animate(fridge, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const handles = el.querySelectorAll('path[d="M9 6h.01"], path[d="M9 14h.01"]');
  animate(handles, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
