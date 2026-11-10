import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector('path[d^="M22 17"]');
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const marks = [el.querySelector('path[d^="M8 14"]'), el.querySelector('path[d^="M14 14"]')];
  animate(marks, { translateY: [-1.5, 0.5, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
