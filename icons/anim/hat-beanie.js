import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const beanie = el.querySelector("g");
  animate(beanie, { translateY: [0, -3, 0], rotate: [0, -10, 4, 0] }, { duration: τmax });

  const pompom = el.querySelector("circle");
  animate(pompom, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
