import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('path[d^="M22 6.675"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const lines = el.querySelectorAll('path[d="M12 10H6"], path[d="M6 14h4"]');
  animate(lines, { opacity: [0, 1], x: [-2, 0] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });

  const seal = el.querySelector("g");
  animate(seal, { scale: [1, 0.85, 1.1, 1] }, { duration: τmid, delay: τmin });
}
