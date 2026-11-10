import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = el.querySelector('path[d^="M17 6"]');
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const latches = el.querySelectorAll('path[d="M8 12v4"], path[d="M16 12v4"]');
  animate(latches, { scaleY: [0, 1.25, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
