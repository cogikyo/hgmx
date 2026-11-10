import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const briefcase = el.querySelector("g");
  animate(briefcase, { translateY: [0, -3, 1, 0] }, { duration: τmax });

  const body = el.querySelector("rect");
  animate(body, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmid, delay: τmin });

  const rollers = el.querySelectorAll('path[d$="20v2"]');
  animate(rollers, { translateX: [0, 2, -0.5, 0] }, { duration: τmid, delay: τmin });
}
