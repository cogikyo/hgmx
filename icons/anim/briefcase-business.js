import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const briefcase = el.querySelector("g");
  animate(briefcase, { translateY: [0, -3, 1, 0] }, { duration: τmax });

  const body = el.querySelector("g g");
  animate(body, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmid, delay: τmin });

  const lid = el.querySelector('path[d^="M22 13"]');
  animate(lid, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
