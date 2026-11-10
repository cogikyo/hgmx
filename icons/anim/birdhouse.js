import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const house = el.querySelector("g");
  animate(house, { scaleY: [1, 0.9, 1.05, 1] }, { duration: τmax });

  const hole = el.querySelector("circle");
  animate(hole, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
