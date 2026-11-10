import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bottle = el.querySelector('path[d^="M9 10c"]');
  animate(bottle, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pump = el.querySelectorAll('path[d="M9 2h7"], path[d="M11 2v4"]');
  animate(pump, { translateY: [0, 1.5, 0] }, { duration: τmin });

  const drops = [el.querySelector("circle"), el.querySelector('path[d="M20 10h.01"]')];
  animate(drops, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
