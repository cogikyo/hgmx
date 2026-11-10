import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const waves = [el.querySelector('path[d^="M13.414"]'), ...el.querySelectorAll("g")];
  animate(waves, { scale: [0.8, 1.1, 1], opacity: [0.5, 1] }, { duration: τmin, delay: stagger(0.03) });

  const slash = el.querySelector('path[d="m2 2 20 20"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
