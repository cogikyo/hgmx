import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const sweep = el.querySelector('path[d^="m13.41"]');
  animate(sweep, { rotate: [0, 360] }, { duration: τmax });

  const blips = [el.querySelector('path[d="M12 18h.01"]'), el.querySelector('path[d="M4 6h.01"]')];
  animate(blips, { scale: [1, 1.75, 1] }, { duration: 0.3, delay: stagger(0.09, { startDelay: τmin }) });
}
