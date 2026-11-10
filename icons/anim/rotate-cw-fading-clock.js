import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const hands = el.querySelector('path[d^="M12 7"]');
  animate(hands, { rotate: [0, 360] }, { duration: τmax });

  const trail = [el.querySelector('path[d^="M21 12"]'), el.querySelector('path[d^="M19 17"]'), el.querySelector('path[d^="M14 20"]')];
  animate(trail, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
