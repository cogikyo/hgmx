import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const card = el.querySelector("g.card");
  animate(card, { translateX: [0, 4, -2, 0] }, { duration: τmax });

  const tick = el.querySelector('path[d^="m16 17"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
