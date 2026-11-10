import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const briefcase = el.querySelector("g");
  animate(briefcase, { translateY: [0, -3, 1, 0] }, { duration: τmax });

  const body = el.querySelector('path[d^="M13.354"]');
  animate(body, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmid, delay: τmin });

  const plus = el.querySelectorAll('path[d="M16 17h6"], path[d="M19 14v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
