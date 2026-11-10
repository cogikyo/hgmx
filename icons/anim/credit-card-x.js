import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const card = el.querySelector("g.card");
  animate(card, { translateX: [0, 4, -2, 0] }, { duration: τmax });

  const x = el.querySelectorAll('path[d^="m21.5"], path[d^="m16.5"]');
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
