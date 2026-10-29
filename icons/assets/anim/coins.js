import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector("g.front");
  animate(front, { translateY: [-6, 1, 0], opacity: [0, 1] }, { duration: τmax });

  const back = el.querySelector("g.back");
  animate(back, { translateY: [0, 1, 0] }, { duration: τmin, delay: τmid });
}
