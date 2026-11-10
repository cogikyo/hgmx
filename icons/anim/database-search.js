import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelectorAll(":scope > path");
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lid = el.querySelector("ellipse");
  animate(lid, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });

  const lens = el.querySelector("g");
  animate(lens, { scale: [1.35, 1], translateX: [0, -4, 2, 0], translateY: [0, -3, -2, 0] }, { duration: τmax });
}
