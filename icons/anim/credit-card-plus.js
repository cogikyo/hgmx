import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const card = el.querySelector("g.card");
  animate(card, { translateX: [0, 4, -2, 0] }, { duration: τmax });

  const plus = el.querySelectorAll('path[d="M16 17h6"], path[d="M19 14v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
