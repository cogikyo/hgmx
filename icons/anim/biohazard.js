import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const sign = el.querySelector("g");
  animate(sign, { rotate: [0, -15, 10, -5, 0] }, { duration: τmax });

  const core = el.querySelector("circle");
  animate(core, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
