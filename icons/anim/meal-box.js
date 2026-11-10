import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = el.querySelectorAll('rect[width="20"], path[d^="M22 12"]');
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const compartments = el.querySelectorAll('rect[width="4"]');
  animate(compartments, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
