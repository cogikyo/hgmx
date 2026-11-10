import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const queen = el.querySelector("g");
  animate(queen, { translateY: [0, -3, 1, 0] }, { duration: τmax });
  const crown = el.querySelector("g.crown");
  animate(crown, { scale: [1, 1.15, 1] }, { duration: τmin, delay: τmid });
}
