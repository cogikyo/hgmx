import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const pumpkin = el.querySelector('path[d^="M17 4"]');
  animate(pumpkin, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const face = el.querySelectorAll('path[d^="M10 11"], path[d^="m18 11"], path[d^="m6 15"]');
  animate(face, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
