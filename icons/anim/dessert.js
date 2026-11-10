import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pudding = el.querySelectorAll("path");
  animate(pudding, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cherry = el.querySelector("circle");
  animate(cherry, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
