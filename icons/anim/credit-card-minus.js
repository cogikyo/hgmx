import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const card = el.querySelector("g.card");
  animate(card, { translateX: [0, 4, -2, 0] }, { duration: τmax });

  const minus = el.querySelector('path[d="M16 17h6"]');
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
