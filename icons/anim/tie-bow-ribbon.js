import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const [bow, tails] = el.querySelectorAll("g");
  animate(bow, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
  animate(tails, { rotate: [0, 8, -6, 0] }, { duration: τmid, delay: τmin });
}
