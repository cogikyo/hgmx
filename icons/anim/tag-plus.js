import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const tag = el.querySelector("g");
  animate(tag, { rotate: [0, -18, 14, -8, 0] }, { duration: τmax });

  const plus = el.querySelectorAll('path[d="M16 13h6"], path[d="M19 10v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
