import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const key = el.querySelector("g");
  animate(key, { translateX: [4, -1, 0], translateY: [-4, 1, 0] }, { duration: τmax });

  const bow = el.querySelector("circle");
  animate(bow, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
